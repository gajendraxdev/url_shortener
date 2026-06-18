package url

import (
	db "url_shortener/config"
)

func CreateShortUrlService(data ShortURLData) (ShortURLData, error) {
	var response ShortURLData
	query := `INSERT INTO urls (original_url, short_code) 
				   VALUES (?, ?) 
				   RETURNING id, original_url, short_code, created_at`
	err := db.DB.QueryRow(query, data.OriginalUrl, data.ShortCode).Scan(&response.ID, &response.OriginalUrl, &response.ShortCode, &response.CreatedAt)

	if err != nil {
		return ShortURLData{}, err
	}

	return response, nil
}

func GetShortUrlService(short_id string) (ShortURLData, error) {
	var result ShortURLData
	err := db.DB.QueryRow("SELECT * FROM urls WHERE short_code = ?", short_id).Scan(&result.ID, &result.OriginalUrl, &result.ShortCode, &result.CreatedAt)

	if err != nil {
		return ShortURLData{}, err
	}

	return result, nil
}
