package url

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	short_id "url_shortener/utils"
)

func CreateShortUrlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data ShortUrlRequest

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request json", http.StatusBadRequest)
		return
	}

	if data.URL == "" {
		http.Error(w, "Provide url is empty or not valid", http.StatusBadRequest)
		return
	}

	var short_code string
	var response ShortURLData
	var err error

	if data.CustomKey != "" {
		short_code = data.CustomKey
		response, err = CreateShortUrlService(ShortURLData{OriginalUrl: data.URL, ShortCode: short_code})
	} else {
		// Retry loop for normal generation
		for i := 0; i < 8; i++ { // max 8 attempts
			code, genErr := short_id.GenerateShortId(5)
			if genErr != nil {
				http.Error(w, "Failed to generate short id", http.StatusInternalServerError)
				return
			}
			short_code = code

			response, err = CreateShortUrlService(ShortURLData{OriginalUrl: data.URL, ShortCode: short_code})
			if err == nil {
				break // success
			}
			// if duplicate, continue to next attempt
		}
	}

	if err != nil {
		http.Error(w, "Failed to create short URL", http.StatusUnprocessableEntity)
		return
	}

	shortUrl := "http://" + r.Host + "/" + response.ShortCode
	log.Println("Short URL", shortUrl)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ShortenResponse{
		ID:          response.ID,
		OriginalUrl: response.OriginalUrl,
		ShortCode:   response.ShortCode,
		ShortUrl:    shortUrl,
		CreatedAt:   response.CreatedAt,
	})
}

func RedirectToOriginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	shortUrlId := strings.TrimPrefix(r.URL.Path, "/")

	if shortUrlId == "" || strings.HasPrefix(shortUrlId, "api") {
		http.NotFound(w, r)
		return
	}

	urlData, err := GetShortUrlService(shortUrlId)

	if err != nil {
		http.Error(w, "Original Url not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, urlData.OriginalUrl, http.StatusFound)
}
