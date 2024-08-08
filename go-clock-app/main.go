package main

import (
	"go-clock-app/dependencies"
	"log"
)

func main() {
	err := dependencies.Initialize()
	if err != nil {
		log.Printf("Failed initialize dependencies: %v", err)
	}
}

