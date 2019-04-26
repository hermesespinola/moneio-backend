package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

// Connect creates a connection to the postgres database
func Connect() *sql.DB {
	fmtStr := "host=%s port=%s user=%s " +
		"password=%s dbname=%s sslmode=disable"
	psqlInfo := fmt.Sprintf(
		fmtStr,
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_DBNAME"),
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

var db *sql.DB

func init() {
	db = Connect()
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}
