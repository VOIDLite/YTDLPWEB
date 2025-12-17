package model

// VideoInfo merepresentasikan informasi dasar video
type VideoInfo struct {
	ID          string                 `json:"id"`
	Title       string                 `json:"title"`
	Duration    float64                `json:"duration"`
	Formats     []FormatInfo           `json:"formats"`
	Thumbnails  []ThumbnailInfo        `json:"thumbnails"`
	RawResponse map[string]interface{} `json:"raw_response"`
}

// FormatInfo merepresentasikan informasi format video/audio
type FormatInfo struct {
	FormatID     string `json:"format_id"`
	FormatNote   string `json:"format_note"`
	Ext          string `json:"ext"`
	Resolution   string `json:"resolution"`
	Quality      int    `json:"quality"`
	FileSize     int64  `json:"filesize,omitempty"`
	AudioOnly    bool   `json:"audio_only"`
	VideoOnly    bool   `json:"video_only"`
	AverageBitrate int  `json:"average_bitrate,omitempty"`
	AudioBitrate int   `json:"audio_bitrate,omitempty"`
	VideoBitrate int   `json:"video_bitrate,omitempty"`
	Height       int   `json:"height,omitempty"`
	Width        int   `json:"width,omitempty"`
	HasVideo     bool  `json:"has_video"`
	HasAudio     bool  `json:"has_audio"`
}

// ThumbnailInfo merepresentasikan informasi thumbnail
type ThumbnailInfo struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// DownloadRequest merepresentasikan permintaan download
type DownloadRequest struct {
	URL          string `json:"url"`
	Format       string `json:"format"` // format_id dari video atau audio
	OutputFolder string `json:"output_folder"`
	Filename     string `json:"filename"`
}

// DownloadResponse merepresentasikan respons download
type DownloadResponse struct {
	Status      string `json:"status"`
	Filename    string `json:"filename,omitempty"`
	Message     string `json:"message,omitempty"`
	Error       string `json:"error,omitempty"`
	DownloadURL string `json:"download_url,omitempty"`
}