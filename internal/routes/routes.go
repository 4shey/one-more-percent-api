package routes

import (
	"net/http"

	"one_more_percent/internal/handlers"
)

func SetupRoutes() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/", handlers.HealthHandler)
	router.HandleFunc("/webhook", handlers.TelegramWebhookHandler)

	return router
}