package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const postsDir = "posts"

type BlogPost struct {
	Date  string
	Title string
	URL   string
}

func extractPostMetadata(filename string) (*BlogPost, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	html := string(content)

	// Extract date from meta tag
	dateRegex := regexp.MustCompile(`<meta\s+name="date"\s+content="([^"]+)"`)
	dateMatch := dateRegex.FindStringSubmatch(html)
	date := ""
	if len(dateMatch) > 1 {
		date = dateMatch[1]
	}

	// Extract title from h1 tag
	titleRegex := regexp.MustCompile(`<h1>([^<]+)</h1>`)
	titleMatch := titleRegex.FindStringSubmatch(html)
	title := ""
	if len(titleMatch) > 1 {
		title = titleMatch[1]
	}

	return &BlogPost{
		Date:  date,
		Title: title,
		URL:   "",
	}, nil
}

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

	var posts []BlogPost
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".html") {
			filename := filepath.Join(postsDir, file.Name())
			post, err := extractPostMetadata(filename)
			if err != nil {
				log.Printf("[ERROR] Failed to extract metadata from %s: %v", file.Name(), err)
				continue
			}
			post.URL = strings.TrimSuffix(file.Name(), ".html")
			posts = append(posts, *post)
		}
	}

	// Sort posts by date, newest first
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date > posts[j].Date
	})

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
