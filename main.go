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
	http.ListenAndServe("0.0.0.0:"+constant.PORT, nil)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("Hello from client side")
	fmt.Fprint(w, "Looks Good :D")
}

func injectRoutes() {
	http.HandleFunc("/", url.RedirectToOriginHandler)
	http.HandleFunc("/api/health", healthCheck)
	http.HandleFunc("/api/shorten", url.CreateShortUrlHandler)
}
