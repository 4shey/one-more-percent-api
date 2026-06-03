package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	database "one_more_percent/internal/db"
	"one_more_percent/internal/routes"
	"one_more_percent/internal/services"
	_ "time/tzdata"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func runMigration() {
	databaseURL := os.Getenv("DATABASE_URL")

	m, err := migrate.New(
		"file:///root/database/migrations",
	databaseURL,
)

	if err != nil {
		log.Fatal("migration init:", err)
	}

	err = m.Up()

	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("migration up:", err)
	}

	fmt.Println("Migration completed ✓")
}

func main() {
	if err := database.Connect(); err != nil {
		log.Fatal("DB connection failed:", err)
	}

	runMigration()

	services.StartScheduler()

	router := routes.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on :%s\n", port)

	log.Fatal(http.ListenAndServe(":"+port, router))
}