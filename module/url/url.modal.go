package url

type ShortURLData struct {
	ID          int    `json:"id"`
	OriginalUrl string `json:"original_url"`
	ShortCode   string `json:"short_code"`
	CreatedAt   string `json:"created_at"`
}

type ShortUrlRequest struct {
	URL       string `json:"url"`
	CustomKey string `json:"custom_key"`
}

type ShortenResponse struct {
	ID          int    `json:"id"`
	OriginalUrl string `json:"original_url"`
	ShortCode   string `json:"short_code"`
	ShortUrl    string `json:"short_url"`
	CreatedAt   string `json:"created_at"`
}
