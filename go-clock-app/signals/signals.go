package signals

import (
	"encoding/json"
	"go-clock-app/clock"
	"net/http"
)

func UpdateSignalHandler(cm *clock.ClockManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}
		var newSignals clock.Signal
		if err := json.NewDecoder(r.Body).Decode(&newSignals); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		cm.UpdateSignals(newSignals)
		w.WriteHeader(http.StatusOK)
	}
}

func UpdateIntervalHandler(cm *clock.ClockManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}
		var newIntervals struct {
			TickInterval string `json:"TickInterval"`
			TockInterval string `json:"TockInterval"`
			BongInterval string `json:"BongInterval"`
		}
		if err := json.NewDecoder(r.Body).Decode(&newIntervals); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		if err := cm.UpdateIntervals(newIntervals.TickInterval, newIntervals.TockInterval, newIntervals.BongInterval); err != nil {
			http.Error(w, "Failed to update intervals: "+err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
