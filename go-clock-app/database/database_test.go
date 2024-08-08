package database

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestNewDatabase(t *testing.T) {
	db := NewDatabase(":memory:")
	defer db.Close()

	var tableName string
	err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='signals'").Scan(&tableName)
	if err != nil {
		t.Fatalf("Failed to verify table creation: %v", err)
	}

	if tableName != "signals" {
		t.Errorf("expected table 'signals', got '%s'", tableName)
	}
}

func TestLogSignal(t *testing.T) {
	db := NewDatabase(":memory:")
	defer db.Close()

	message := "test message"
	db.LogSignal(message)

	var storedMessage string
	err := db.QueryRow("SELECT message FROM signals WHERE message = ?", message).Scan(&storedMessage)
	if err != nil {
		t.Fatalf("Failed to retrieve message from database: %v", err)
	}

	if storedMessage != message {
		t.Errorf("expected '%s', got '%s'", message, storedMessage)
	}
}
func TestLogSignalErrorHandling(t *testing.T) {
	db := NewDatabase(":memory:")
	db.Close()

	message := "test message after close"
	db.LogSignal(message)

	err := db.Ping()
	if err == nil {
		t.Error("expected database to be closed")
	}
}
