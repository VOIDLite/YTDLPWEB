package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"yt-dlp-web/internal/handler"
)

func main() {
	// Port server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

	// Setup routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/info", handler.GetVideoInfoHandler)
	http.HandleFunc("/download", handler.DownloadHandler)
	http.HandleFunc("/downlink", handler.DownlinkHandler)

	// Serve static files
	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve downloads
	downloadFS := http.FileServer(http.Dir("./downloads/"))
	http.Handle("/downloads/", http.StripPrefix("/downloads/", downloadFS))

	// Mulai server
	fmt.Printf("Server berjalan di :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// homeHandler menyajikan halaman utama
func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./templates/index.html")
}