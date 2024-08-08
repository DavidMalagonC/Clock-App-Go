package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	*sql.DB
}

func NewDatabase(dbPath string) *Database {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS signals (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		message TEXT
	)`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	return &Database{DB: db}
}

func (db *Database) LogSignal(message string) {
	_, err := db.Exec("INSERT INTO signals (message) VALUES (?)", message)
	if err != nil {
		log.Printf("Failed to insert signal: %v", err)
	}
}
