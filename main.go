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
	log.Printf("[INFO] %s %s - Remote: %s, User-Agent: %s",
		r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())

	files, err := os.ReadDir(postsDir)
	if err != nil {
		log.Printf("[ERROR] Failed to read posts directory: %v", err)
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

	log.Printf("[INFO] Found %d blog posts", len(posts))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := ListTmpl.Execute(w, posts); err != nil {
		log.Printf("[ERROR] Template execution failed: %v", err)
		return
	}

	log.Printf("[INFO] Successfully rendered blog list page")
}

func serveBlogPost(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] %s %s - Remote: %s, User-Agent: %s",
		r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())

	title := strings.TrimPrefix(r.URL.Path, "/blog/")

	if title == "" {
		log.Printf("[ERROR] Failed to serve blog post.")
		http.Redirect(w, r, "/blog", http.StatusMovedPermanently)
		return
	}

	filename := filepath.Join(postsDir, title+".html")

	content, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("[ERROR] Failed to load content for blog %s", title)
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

	log.Printf("[INFO] Succesfully loaded content for blog post")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := PostTmpl.Execute(w, data); err != nil {
		log.Printf("[ERROR] Template execution failed: %v", err)
		return
	}
}

func redirectToBlog(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/blog", http.StatusMovedPermanently)
		return
	}

	http.NotFound(w, r)
}
