package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

// DownlinkHandler menangani permintaan download file yang telah diunduh
func DownlinkHandler(w http.ResponseWriter, r *http.Request) {
	// Dapatkan parameter reff dari query
	reff := r.URL.Query().Get("reff")
	if reff == "" {
		http.Error(w, "Parameter reff tidak ditemukan", http.StatusBadRequest)
		return
	}

	// Decode URL parameter untuk mencegah path traversal
	decodedReff, err := url.QueryUnescape(reff)
	if err != nil {
		http.Error(w, "Parameter reff tidak valid", http.StatusBadRequest)
		return
	}

	// Validasi nama file untuk mencegah path traversal
	if filepath.Base(decodedReff) != decodedReff {
		http.Error(w, "Nama file tidak valid", http.StatusBadRequest)
		return
	}

	// Buat path lengkap ke file download
	filePath := filepath.Join("./downloads", decodedReff)

	// Cek apakah file ada
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File tidak ditemukan", http.StatusNotFound)
		return
	}

	// Set header untuk download
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", decodedReff))
	w.Header().Set("Content-Type", "application/octet-stream")

	// Kirim file ke client
	http.ServeFile(w, r, filePath)

	// Hapus file setelah dikirim (opsional, bisa dihapus setelah download selesai)
	// Kita bisa hapus file setelah response selesai dikirim
	go func() {
		// Tunggu sebentar untuk memastikan response telah dikirim
		// Note: Dalam implementasi production, sebaiknya menggunakan solusi yang lebih robust
		// seperti menandai file untuk dihapus nanti, atau menggunakan mekanisme async
		time.Sleep(5 * time.Second)
		if err := os.Remove(filePath); err != nil {
			fmt.Printf("Gagal menghapus file %s: %v\n", filePath, err)
		} else {
			fmt.Printf("File %s berhasil dihapus setelah download\n", filePath)
		}
	}()
}