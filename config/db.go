package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
	"url_shortener/constant"
)

var DB *sql.DB

func InitializeDatabase() {

	var err error

	DB, err = sql.Open("sqlite", constant.Database)
	if err != nil {
		log.Fatalln("Failed to open database connection", err)
	}

	DB.SetConnMaxLifetime(time.Minute * 5)
	DB.SetMaxOpenConns(1)
	DB.SetMaxIdleConns(1)

	// === Important for concurrency ===
	DB.Exec("PRAGMA journal_mode = WAL;")
	DB.Exec("PRAGMA busy_timeout = 5000;") // wait up to 5 seconds on lock
	// DB.Exec("CREATE INDEX IF NOT EXISTS idx_short_code ON urls(short_code);")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := DB.PingContext(ctx); err != nil {
		log.Fatalln("Failed to connect database", err)
		return
	}

	fmt.Println("Database connection established")

	schema := `
	CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		original_url TEXT NOT NULL,
		short_code TEXT NOT NULL UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_short_code ON urls(short_code);
	`

	if _, err := DB.ExecContext(ctx, schema); err != nil {
		log.Fatalln("Failed to initialize tables", err)
	}
}
