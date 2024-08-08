package dependencies

import (
	"fmt"
	"go-clock-app/clock"
	"go-clock-app/database"
	"go-clock-app/signals"
	"net/http"
	"os"
)

func Initialize() error {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./signals.db"
	}
	db := database.NewDatabase(dbPath)
	cm := clock.NewClockManager(db)

	go cm.Run()

	http.HandleFunc("/update-signals", signals.UpdateSignalHandler(cm))
	http.HandleFunc("/update-intervals", signals.UpdateIntervalHandler(cm))
	fmt.Println("Server started on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
		return err
	}
	return nil
}
