package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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

	// Railway will inject PORT automatically
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on :%s\n", port)

	log.Fatal(http.ListenAndServe(":"+port, router))
}