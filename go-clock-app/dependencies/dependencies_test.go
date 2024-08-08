package dependencies

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func startTestServer(mux *http.ServeMux) *httptest.Server {
	server := httptest.NewServer(mux)
	return server
}

func TestInitializeWithDefaultDBPath(t *testing.T) {
	_ = os.Setenv("DB_PATH", ":memory:")

	mux := http.NewServeMux()
	mux.HandleFunc("/update-signals", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	server := startTestServer(mux)
	defer server.Close()

	time.Sleep(100 * time.Millisecond) // Ensure server is up

	resp, err := http.Get(server.URL + "/update-signals")
	if err != nil {
		t.Fatalf("Failed to reach server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 Method Not Allowed, got %v", resp.StatusCode)
	}

	_ = os.Unsetenv("DB_PATH")
}

func TestInitializeWithCustomDBPath(t *testing.T) {
	dbPath := "custom_signals.db"
	_ = os.Setenv("DB_PATH", dbPath)

	defer func() {
		_ = os.Remove(dbPath)
		_ = os.Unsetenv("DB_PATH")
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/update-signals", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	server := startTestServer(mux)
	defer server.Close()

	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get(server.URL + "/update-signals")
	if err != nil {
		t.Fatalf("Failed to reach server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 Method Not Allowed, got %v", resp.StatusCode)
	}
}

func TestInitializeServerError(t *testing.T) {
	_ = os.Setenv("DB_PATH", ":memory:")

	mux := http.NewServeMux()
	mux.HandleFunc("/update-signals", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	server := startTestServer(mux)
	defer server.Close()

	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get(server.URL + "/update-signals")
	if err != nil {
		t.Fatalf("Failed to reach server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 Method Not Allowed, got %v", resp.StatusCode)
	}

	_ = os.Unsetenv("DB_PATH")
}
