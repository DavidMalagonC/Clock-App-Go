package clock

import (
	"go-clock-app/database"
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {
	db := database.NewDatabase(":memory:")

	cm := NewManager(db)

	if cm.signals.TickMessage != "tick" {
		t.Errorf("expected 'tick', got '%s'", cm.signals.TickMessage)
	}

	if cm.signals.TockMessage != "tock" {
		t.Errorf("expected 'tock', got '%s'", cm.signals.TockMessage)
	}

	if cm.signals.BongMessage != "bong" {
		t.Errorf("expected 'bong', got '%s'", cm.signals.BongMessage)
	}

	if cm.GetTickInterval() != 1*time.Second {
		t.Errorf("expected 1s, got %s", cm.GetTickInterval())
	}

	if cm.GetTockInterval() != 1*time.Minute {
		t.Errorf("expected 1m, got %s", cm.GetTockInterval())
	}

	if cm.GetBongInterval() != 1*time.Hour {
		t.Errorf("expected 1h, got %s", cm.GetBongInterval())
	}
}

func TestUpdateSignals(t *testing.T) {
	db := database.NewDatabase(":memory:")
	cm := NewManager(db)
	newSignals := Signal{
		TickMessage: "quack",
		TockMessage: "dong",
		BongMessage: "boom",
	}

	cm.UpdateSignals(newSignals)

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
		t.Error("UpdateSignals did not send updated signals on time")
	}
}

func TestUpdateIntervals(t *testing.T) {
	db := database.NewDatabase(":memory:")
	cm := NewManager(db)

	err := cm.UpdateIntervals("5s", "3m", "2h")
	if err != nil {
		t.Fatalf("Failed to update intervals: %v", err)
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

func TestLogSignal(t *testing.T) {
	db := database.NewDatabase(":memory:")
	cm := NewManager(db)

	message := "test message"
	cm.LogSignal(message)

	row := db.QueryRow("SELECT message FROM signals WHERE message = ?", message)
	var storedMessage string
	err := row.Scan(&storedMessage)
	if err != nil {
		t.Fatalf("Failed to retrieve message from database: %v", err)
	}

	if storedMessage != message {
		t.Errorf("expected '%s', got '%s'", message, storedMessage)
	}
}
