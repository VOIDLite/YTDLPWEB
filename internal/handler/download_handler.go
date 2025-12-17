package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"yt-dlp-web/internal/model"
	"yt-dlp-web/internal/utils"
)

// DownloadHandler menangani permintaan download video/audio
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method tidak diperbolehkan", http.StatusMethodNotAllowed)
		return
	}

	var req model.DownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Gagal membaca permintaan", http.StatusBadRequest)
		return
	}

	// Validasi URL
	if req.URL == "" {
		http.Error(w, "URL tidak boleh kosong", http.StatusBadRequest)
		return
	}

	if !utils.IsValidURL(req.URL) {
		http.Error(w, "URL tidak valid", http.StatusBadRequest)
		return
	}

	// Set header untuk SSE (Server Sent Events) untuk log real-time
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Buat saluran untuk mengirim log
	done := make(chan bool)

	// Fungsi untuk mengirim log ke klien
	sendLog := func(message string) {
		fmt.Fprintf(w, "data: %s\n\n", message)
		flusher, ok := w.(http.Flusher)
		if ok {
			flusher.Flush()
		}
	}

	// Kirim log awal
	sendLog("Memulai persiapan download...")

	// Buat nama file jika tidak disediakan
	if req.Filename == "" {
		req.Filename = "%(title)s.%(ext)s"
	}

	// Ekstensi file untuk nama output
	outputName := req.Filename

	// Direktori output
	outputDir := req.OutputFolder
	if outputDir == "" {
		outputDir = "./downloads"
	}

	// Pastikan direktori output ada
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		sendLog(fmt.Sprintf("Error membuat direktori output: %v", err))
		return
	}

	// Jalankan download dalam goroutine agar bisa streaming log
	go func() {
		defer close(done)

		// Bangun perintah yt-dlp
		args := []string{
			"-f", req.Format,
			"-o", filepath.Join(outputDir, outputName),
			"--newline",
		}

		// Tambahkan URL
		args = append(args, req.URL)

		// Buat perintah
		cmd := exec.Command("yt-dlp", args...)

		// Tangani output dan error
		stdout, _ := cmd.StdoutPipe()
		stderr, _ := cmd.StderrPipe()

		sendLog(fmt.Sprintf("Menjalankan perintah: yt-dlp %s", strings.Join(args, " ")))

		err := cmd.Start()
		if err != nil {
			sendLog(fmt.Sprintf("Error memulai perintah: %v", err))
			return
		}

		// Gabungkan stdout dan stderr
		reader := io.MultiReader(stdout, stderr)

		// Buffer untuk membaca output
		buf := make([]byte, 1024)

		// Baca output saat sedang berlangsung
		for {
			n, err := reader.Read(buf)
			if n > 0 {
				// Kirim log ke client
				logLine := string(buf[:n])
				lines := strings.Split(logLine, "\n")

				for _, line := range lines {
					line = strings.TrimSpace(line)
					if line != "" {
						sendLog(line)
					}
				}
			}

			if err != nil {
				if err == io.EOF {
					break
				}
				sendLog(fmt.Sprintf("Error membaca output: %v", err))
				break
			}
		}

		// Tunggu perintah selesai
		err = cmd.Wait()
		if err != nil {
			sendLog(fmt.Sprintf("Perintah selesai dengan error: %v", err))
		} else {
			sendLog("Proses download selesai!")
		}
	}()

	// Tunggu hingga selesai atau timeout
	select {
	case <-done:
		// Download selesai
		sendLog("FINISHED")
	case <-time.After(30 * time.Minute):
		// Timeout setelah 30 menit
		sendLog("TIMEOUT: Proses download terlalu lama")
	}
}