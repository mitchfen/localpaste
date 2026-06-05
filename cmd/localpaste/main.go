package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"localpaste/internal/paste"
)

type Page struct {
	Paste *paste.Paste
}

func init() {
	paste.StartCleanupRoutine()
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	page := &Page{}

	if r.Method == http.MethodPost {
		content := r.FormValue("content")
		if content == "" {
			http.Error(w, "Content cannot be empty", http.StatusBadRequest)
			return
		}

		m := paste.GetMutex()
		m.Lock()
		p, err := paste.Write(content)
		m.Unlock()

		if err != nil {
			log.Printf("Write error: %v", err)
			http.Error(w, "Failed to save paste", http.StatusInternalServerError)
			return
		}

		page.Paste = p

		// Return just the paste section HTML for HTMX
		w.Header().Set("Content-Type", "text/html")
		tmpl, err := template.New("paste").Parse(`
<div class="paste-section">
<div class="paste-header">
<p>Current paste (expires in {{.ExpiresIn}})</p>
</div>
<div class="paste-content" id="paste-text">{{.Content}}</div>
</div>
`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, page.Paste)
		return
	}

	m := paste.GetMutex()
	m.Lock()
	p, err := paste.Read()
	m.Unlock()

	if err != nil {
		log.Printf("Read error: %v", err)
		http.Error(w, "Error retrieving paste", http.StatusInternalServerError)
		return
	}
	if p != nil {
		page.Paste = p
	}

	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, page)
}

func main() {
	http.HandleFunc("/", handleIndex)
	http.Handle("/manifest.json", http.FileServer(http.Dir("web/templates")))
	http.Handle("/android-chrome-192x192.png", http.FileServer(http.Dir(".")))
	http.Handle("/android-chrome-512x512.png", http.FileServer(http.Dir(".")))
	http.Handle("/screenshot.png", http.FileServer(http.Dir(".")))

	fmt.Println("LocalPaste running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
