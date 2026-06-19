package main

import (
	"fmt"
	"log"
	"net/http"
	db "url_shortener/config"
	"url_shortener/constant"
	"url_shortener/module/url"

	_ "modernc.org/sqlite"
)

func main() {
	db.InitializeDatabase()
	injectRoutes()

	log.Println("Server is running on PORT", constant.PORT)
	http.ListenAndServe(constant.PORT, nil)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("Hello from client side")
	fmt.Fprint(w, "Looks Good :D")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// Serve index at root
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "web/index.html")
		return
	}

	// For other non-api paths, delegate to redirect handler
	url.RedirectToOriginHandler(w, r)
}

func injectRoutes() {
	// Serve static assets
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Root handler serves UI and delegates redirects
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/api/health", healthCheck)
	http.HandleFunc("/api/shorten", url.CreateShortUrlHandler)
}
