package main

import (
	"fmt"
	"log"
	"net/http"

	database "one_more_percent/internal/db"
	"one_more_percent/internal/routes"
	"one_more_percent/internal/services"
)

func main() {
	// Connect to PostgreSQL.
	if err := database.Connect(); err != nil {
		log.Fatal("DB connection failed:", err)
	}

	// Start background scheduler (reminders + midnight recap).
	services.StartScheduler()

	router := routes.SetupRoutes()

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", router)
}