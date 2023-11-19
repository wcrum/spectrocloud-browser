package main

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Response(w http.ResponseWriter, r *http.Request, args ServerArgs) {
	// Check if file exists and is not a directory
	path := filepath.Clean(r.URL.Path)
	if fileInfo, err := os.Stat(args.Directory + path); err != nil {
		http.NotFound(w, r)
		return
	} else if fileInfo.IsDir() {
		files, err := os.ReadDir(args.Directory + path)
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
			Path:  strings.TrimRight(path, "/"),
			Files: fileNames,
		})
		return
	}

	// Set the header and write the file content
	filename := strings.Split(path, "/")[0]

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/octet-stream")
	file, err := os.Open(args.Directory + path)
	if err != nil {
		http.Error(w, "File not found.", 404)
		return
	}
	defer file.Close()
	io.Copy(w, file)
}

func main() {

	args := NewServerArgs()

	if args.Directory == "basic" {
		http.Handle("/", BasicAuthMiddleware(args, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Response(w, r, args)
		})))
	} else {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			Response(w, r, args)
		})
	}

	if args.TLSCertFile != "" && args.TLSKeyFile != "" {

		http.ListenAndServeTLS(":"+args.Port, args.TLSCertFile, args.TLSKeyFile, nil)
	}

	http.ListenAndServe(":"+args.Port, nil)
}
