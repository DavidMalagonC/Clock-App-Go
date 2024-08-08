package signals

import (
	"bytes"
	"encoding/json"
	"go-clock-app/clock"
	"go-clock-app/database"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUpdateSignalHandler(t *testing.T) {
	db := database.NewDatabase(":memory:")
	cm := clock.NewManager(db)
	handler := UpdateSignalHandler(cm)

	newSignals := clock.Signal{
		TickMessage: "quack",
		TockMessage: "dong",
		BongMessage: "boom",
	}

	jsonData, err := json.Marshal(newSignals)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest("POST", "/update-signals", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	select {
	case updatedSignals := <-cm.Updates:
		if updatedSignals.TickMessage != "quack" {
			t.Errorf("expected 'quack', got '%s'", updatedSignals.TickMessage)
		}

		if updatedSignals.TockMessage != "dong" {
			t.Errorf("expected 'dong', got '%s'", updatedSignals.TockMessage)
		}

		if updatedSignals.BongMessage != "boom" {
			t.Errorf("expected 'boom', got '%s'", updatedSignals.BongMessage)
		}
	case <-time.After(1 * time.Second):
		t.Error("UpdateSignalHandler did not send updated signals on time")
	}
}

func TestUpdateIntervalHandler(t *testing.T) {
	db := database.NewDatabase(":memory:")
	cm := clock.NewManager(db)
	handler := UpdateIntervalHandler(cm)

	newIntervals := struct {
		TickInterval string `json:"TickInterval"`
		TockInterval string `json:"TockInterval"`
		BongInterval string `json:"BongInterval"`
	}{
		TickInterval: "5s",
		TockInterval: "3m",
		BongInterval: "2h",
	}

	jsonData, err := json.Marshal(newIntervals)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest("POST", "/update-intervals", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if cm.GetTickInterval() != 5*time.Second {
		t.Errorf("expected 5s, got %s", cm.GetTickInterval())
	}

	if cm.GetTockInterval() != 3*time.Minute {
		t.Errorf("expected 3m, got %s", cm.GetTockInterval())
	}

	if cm.GetBongInterval() != 2*time.Hour {
		t.Errorf("expected 2h, got %s", cm.GetBongInterval())
	}
}
