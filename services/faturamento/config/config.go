package config

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const defaultEstoqueAPIURL = "http://localhost:8080"

func ConnectDB() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func GetEstoqueAPIURL() string {
	apiURL := os.Getenv("ESTOQUE_API_URL")
	if apiURL == "" {
		return defaultEstoqueAPIURL
	}
	return apiURL
}
