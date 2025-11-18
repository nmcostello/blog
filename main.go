package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const postsDir = "posts"

func main() {
	mux := http.NewServeMux()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux.HandleFunc("/blog", listBlogPosts)
	mux.HandleFunc("/blog/", serveBlogPost)
	mux.Handle("/pictures/", http.StripPrefix("/pictures/", http.FileServer(http.Dir("pictures"))))
	mux.HandleFunc("/", redirectToBlog)

	fmt.Printf("Starting blog server on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func listBlogPosts(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(postsDir)
	if err != nil {
		http.Error(w, "Failed to read posts directory", http.StatusInternalServerError)
		return
	}

	var posts []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".html") {
			url := strings.TrimSuffix(file.Name(), ".html")
			posts = append(posts, url)
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	ListTmpl.Execute(w, posts)
}

func serveBlogPost(w http.ResponseWriter, r *http.Request) {
	title := strings.TrimPrefix(r.URL.Path, "/blog/")

	if title == "" {
		http.Redirect(w, r, "/blog", http.StatusMovedPermanently)
		return
	}

	filename := filepath.Join(postsDir, title+".html")

	content, err := os.ReadFile(filename)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	data := struct {
		Title   string
		Content template.HTML
	}{
		Title:   title,
		Content: template.HTML(content),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	PostTmpl.Execute(w, data)
}

func redirectToBlog(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/blog", http.StatusMovedPermanently)
		return
	}

	http.NotFound(w, r)
}
