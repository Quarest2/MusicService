package db

import (
	"database/sql"
	_ "embed"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var db *sql.DB = nil

//go:embed init.sql
var initSQL string

func Connect() {
	log.Printf("Connecting to the database")

	username, ok1 := os.LookupEnv("DB_USERNAME")
	password, ok2 := os.LookupEnv("DB_PASSWORD")
	databaseIP, ok3 := os.LookupEnv("DB_HOST")
	databaseName, ok4 := os.LookupEnv("DB_NAME")

	if !ok1 || !ok2 || !ok3 || !ok4 {
		log.Printf("Database environment variables are not set")
		return
	}

	connStr := "user=" + username + " password=" + password + " dbname=" + databaseName + " host=" + databaseIP + " sslmode=disable"

	var err error

	db, err = sql.Open("postgres", connStr)

	ErrorHandler(err, "Error while connecting to the database")

	setupDB()
}

func setupDB() {
	log.Printf("Creating the table")

	if db == nil {
		log.Printf("Database is not connected")
		return
	}

	_, err := db.Exec(initSQL)

	ErrorHandler(err, "Error while creating the table")
}
