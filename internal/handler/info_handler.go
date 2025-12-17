package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"yt-dlp-web/internal/utils"
)

// GetVideoInfoHandler menangani permintaan informasi video
func GetVideoInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method tidak diperbolehkan", http.StatusMethodNotAllowed)
		return
	}

	// Ambil URL dari body permintaan
	url := r.FormValue("url")
	if url == "" {
		http.Error(w, "URL tidak boleh kosong", http.StatusBadRequest)
		return
	}

	// Validasi bahwa URL adalah URL YouTube valid
	if !utils.IsValidURL(url) {
		http.Error(w, "URL tidak valid", http.StatusBadRequest)
		return
	}

	// Panggil yt-dlp untuk mendapatkan informasi video
	videoInfo, err := utils.GetVideoInfoFromYtDlp(url)
	if err != nil {
		log.Printf("Error mendapatkan informasi video: %v", err)
		http.Error(w, fmt.Sprintf("Gagal mendapatkan informasi video: %v", err), http.StatusInternalServerError)
		return
	}

	// Set header JSON dan kirim respons
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videoInfo)
}

