package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Retrieve environment variables
	dir := os.Getenv("FILES_DIR")
	if dir == "" {
		dir = "./public" // Default value if not set
	}

	// Download file handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check if file exists and is not a directory
		if fileInfo, err := os.Stat(dir + r.URL.Path); err != nil {
			http.NotFound(w, r)
			return
		} else if fileInfo.IsDir() {
			files, err := os.ReadDir(dir + r.URL.Path)
			if err != nil {
				http.Error(w, "Error reading directory", http.StatusInternalServerError)
				return
			}

			var fileNames []string
			for _, file := range files {
				fileNames = append(fileNames, file.Name())
			}

			tmpl, err := template.New("dirlist").Parse(dirListTemplate)
			if err != nil {
				http.Error(w, "Error rendering template", http.StatusInternalServerError)
				return
			}

			tmpl.Execute(w, struct {
				Path  string
				Files []string
			}{
				Path:  strings.TrimRight(r.URL.Path, "/"),
				Files: fileNames,
			})
			return
		}

		// Set the header and write the file content
		filename := strings.Split(r.URL.Path, "/")[0]

		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		w.Header().Set("Content-Type", "application/octet-stream")
		file, err := os.Open(dir + r.URL.Path)
		if err != nil {
			http.Error(w, "File not found.", 404)
			return
		}
		defer file.Close()
		io.Copy(w, file)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not set
	}

	certFile := os.Getenv("TLS_CERT_FILE")
	keyFile := os.Getenv("TLS_KEY_FILE")
	fmt.Println("test")
	if certFile != "" || keyFile != "" {
		fmt.Println("running on ssl")
		http.ListenAndServeTLS(":"+port, certFile, keyFile, nil)
	}

	http.ListenAndServe(":"+port, nil)
}
