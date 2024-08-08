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
	cm := clock.NewManager(db)

	go cm.Run()

	http.HandleFunc("/update-signals", signals.UpdateSignalHandler(cm))
	http.HandleFunc("/update-intervals", signals.UpdateIntervalHandler(cm))
	
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
		return err
	}
	fmt.Println("Server started on port 8080")
	return nil
}
