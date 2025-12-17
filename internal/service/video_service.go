package service

import (
	"fmt"
	"os/exec"
	"strings"
	"yt-dlp-web/internal/model"
)

// VideoService menyediakan fungsionalitas untuk mengelola video
type VideoService struct{}

// GetVideoInfo mendapatkan informasi video dari URL
func (s *VideoService) GetVideoInfo(url string) (*model.VideoInfo, error) {
	// Validasi bahwa URL adalah URL YouTube valid
	if !isValidURL(url) {
		return nil, fmt.Errorf("URL tidak valid")
	}

	// Panggil yt-dlp untuk mendapatkan informasi video
	videoInfo, err := s.getVideoInfoFromYtDlp(url)
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan informasi video: %v", err)
	}

	return videoInfo, nil
}

// DownloadVideo men-download video dengan format tertentu
func (s *VideoService) DownloadVideo(req model.DownloadRequest) error {
	// Validasi URL
	if req.URL == "" {
		return fmt.Errorf("URL tidak boleh kosong")
	}

	if !isValidURL(req.URL) {
		return fmt.Errorf("URL tidak valid")
	}

	// Bangun perintah yt-dlp
	args := []string{
		"-f", req.Format,
		"-o", req.Filename,
		"--newline",
	}

	// Tambahkan URL
	args = append(args, req.URL)

	// Buat dan jalankan perintah
	cmd := exec.Command("yt-dlp", args...)
	
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("gagal men-download video: %v", err)
	}

	return nil
}

// isValidURL memvalidasi apakah URL adalah URL YouTube valid
func isValidURL(url string) bool {
	// Validasi dasar URL YouTube
	return strings.Contains(url, "youtube.com/watch") || 
		   strings.Contains(url, "youtu.be/") ||
		   strings.Contains(url, "youtube.com/shorts/")
}

// getVideoInfoFromYtDlp mendapatkan informasi video dari yt-dlp
func (s *VideoService) getVideoInfoFromYtDlp(url string) (*model.VideoInfo, error) {
	// Buat perintah yt-dlp untuk mendapatkan info video
	cmd := exec.Command("yt-dlp", "-j", "--no-playlist", url)
	
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("gagal menjalankan yt-dlp: %v", err)
	}

	// Di sini sebenarnya kita perlu menguraikan JSON output ke model.VideoInfo
	// Tapi karena struktur datanya kompleks, kita akan gunakan fungsi dari handler
	// atau buat fungsi parsing yang sama di sini
	
	// Untuk saat ini, kita lempar error karena implementasi parsing
	// sebaiknya dilakukan di satu tempat saja (misalnya di handler)
	return nil, fmt.Errorf("fungsi parsing perlu diimplementasikan di sini atau diambil dari handler")
}