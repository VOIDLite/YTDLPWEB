package utils

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"yt-dlp-web/internal/model"
)

// GetVideoInfoFromYtDlp mendapatkan informasi video dari yt-dlp
func GetVideoInfoFromYtDlp(url string) (*model.VideoInfo, error) {
	// Buat perintah yt-dlp untuk mendapatkan info video
	cmd := exec.Command("yt-dlp", "-j", "--no-playlist", url)
	
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("gagal menjalankan yt-dlp: %v", err)
	}

	// Parse output JSON
	var rawData map[string]interface{}
	err = json.Unmarshal(output, &rawData)
	if err != nil {
		return nil, fmt.Errorf("gagal menguraikan JSON: %v", err)
	}

	// Konversi ke struct VideoInfo
	videoInfo := &model.VideoInfo{
		ID:          getStringValue(rawData, "id"),
		Title:       getStringValue(rawData, "title"),
		Duration:    getFloatValue(rawData, "duration"),
		Thumbnails:  parseThumbnails(rawData),
		RawResponse: rawData,
	}

	// Parsing format info
	formats, ok := rawData["formats"].([]interface{})
	if ok {
		videoInfo.Formats = parseFormats(formats)
	} else if format, ok := rawData["format"].([]interface{}); ok {
		videoInfo.Formats = parseFormats(format)
	}

	return videoInfo, nil
}

// DownloadVideoWithProgress men-download video dengan format tertentu dan callback untuk progress
func DownloadVideoWithProgress(req model.DownloadRequest, logCallback func(string)) error {
	// Validasi URL
	if req.URL == "" {
		return fmt.Errorf("URL tidak boleh kosong")
	}

	// Bangun perintah yt-dlp
	args := []string{
		"-f", req.Format,
		"-o", req.Filename,
		"--newline",
	}

	// Tambahkan URL
	args = append(args, req.URL)

	// Buat perintah
	cmd := exec.Command("yt-dlp", args...)

	// Tangani output dan error
	_, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("gagal membuat pipe untuk stdout: %v", err)
	}
	_, err = cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("gagal membuat pipe untuk stderr: %v", err)
	}

	// Kirim log awal
	if logCallback != nil {
		logCallback(fmt.Sprintf("Menjalankan perintah: yt-dlp %s", strings.Join(args, " ")))
	}

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("gagal memulai perintah: %v", err)
	}

	// Gabungkan stdout dan stderr
	// TODO: Implementasi streaming log

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("perintah selesai dengan error: %v", err)
	}

	return nil
}

// isValidURL memvalidasi apakah URL adalah URL YouTube valid
func IsValidURL(url string) bool {
	// Validasi dasar URL YouTube
	return strings.Contains(url, "youtube.com/watch") || 
		   strings.Contains(url, "youtu.be/") ||
		   strings.Contains(url, "youtube.com/shorts/")
}

// getStringValue membantu mendapatkan nilai string dari map interface
func getStringValue(data map[string]interface{}, key string) string {
	if value, ok := data[key]; ok {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

// getFloatValue membantu mendapatkan nilai float dari map interface
func getFloatValue(data map[string]interface{}, key string) float64 {
	if value, ok := data[key]; ok {
		if f, ok := value.(float64); ok {
			return f
		}
	}
	return 0.0
}

// parseThumbnails menguraikan informasi thumbnail dari data raw
func parseThumbnails(data map[string]interface{}) []model.ThumbnailInfo {
	var thumbnails []model.ThumbnailInfo
	
	if thumbnailData, ok := data["thumbnails"].([]interface{}); ok {
		for _, thumbItem := range thumbnailData {
			if thumbMap, ok := thumbItem.(map[string]interface{}); ok {
				thumbnail := model.ThumbnailInfo{
					URL:    getStringValue(thumbMap, "url"),
					Width:  getIntValue(thumbMap, "width"),
					Height: getIntValue(thumbMap, "height"),
				}
				thumbnails = append(thumbnails, thumbnail)
			}
		}
	}
	
	return thumbnails
}

// getIntValue membantu mendapatkan nilai integer dari map interface
func getIntValue(data map[string]interface{}, key string) int {
	if value, ok := data[key]; ok {
		if f, ok := value.(float64); ok {
			return int(f)
		}
	}
	return 0
}

// parseFormats menguraikan format video/audio
func parseFormats(formatsData []interface{}) []model.FormatInfo {
	var formats []model.FormatInfo
	
	for _, formatItem := range formatsData {
		if formatMap, ok := formatItem.(map[string]interface{}); ok {
			format := model.FormatInfo{
				FormatID:     getStringValue(formatMap, "format_id"),
				FormatNote:   getStringValue(formatMap, "format_note"),
				Ext:          getStringValue(formatMap, "ext"),
				Resolution:   getStringValue(formatMap, "resolution"),
				Quality:      getIntValue(formatMap, "quality"),
				AudioBitrate: getIntValue(formatMap, "abr"),
				AverageBitrate: getIntValue(formatMap, "average_bitrate"),
				VideoBitrate: getIntValue(formatMap, "vbr"),
				Height:       getIntValue(formatMap, "height"),
				Width:        getIntValue(formatMap, "width"),
			}
			
			if fileSize, ok := formatMap["filesize"]; ok {
				if f, ok := fileSize.(float64); ok {
					format.FileSize = int64(f)
				}
			}
			
			// Deteksi apakah format memiliki video atau audio
			videoCodec := getStringValue(formatMap, "vcodec")
			audioCodec := getStringValue(formatMap, "acodec")
			
			format.HasVideo = videoCodec != "none"
			format.HasAudio = audioCodec != "none"
			
			// Deteksi apakah format hanya audio atau hanya video
			format.AudioOnly = !format.HasVideo && format.HasAudio
			format.VideoOnly = format.HasVideo && !format.HasAudio
			
			formats = append(formats, format)
		}
	}
	
	return formats
}