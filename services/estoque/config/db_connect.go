package config

import(
	"database/sql"
	"os"
	"log"

    _ "github.com/lib/pq"
)

var db *sql.DB
var err error

func Db_connect() *sql.DB {
	connStr := os.Getenv("DB_URL")
	
	if connStr == ""{
		log.Fatal("DB_URL env not configured")
	}
	
	db, err = sql.Open("postgres", connStr)
	
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	
	return db
}