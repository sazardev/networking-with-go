// HTTP File Upload/Download Example
// Serves a multipart upload endpoint and a download endpoint that
// delegates to http.ServeFile for range-request support.
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// 32 MB is the memory ceiling for parsed form parts; larger files
	// spill over to temporary disk files automatically.
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "bad upload", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "missing file field", http.StatusBadRequest)
		return
	}
	defer file.Close()

	dst, err := os.Create("uploads/" + filepath.Base(header.Filename))
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	written, err := io.Copy(dst, file)
	if err != nil {
		http.Error(w, "write failed", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "uploaded %d bytes as %s\n", written, header.Filename)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	name := filepath.Base(r.URL.Query().Get("name"))
	http.ServeFile(w, r, filepath.Join("uploads", name))
}

func main() {
	if err := os.MkdirAll("uploads", 0o755); err != nil {
		log.Fatalf("could not create uploads dir: %v", err)
	}
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download", downloadHandler)
	fmt.Println("listening on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
