package main

import (
	"fmt"
	"log"
	"net/http"

	"example.com/m/v2/src/db"
	"example.com/m/v2/src/events"
	"example.com/m/v2/src/server/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	// Initialize Redis
	db.InitRedis()

	// Start Redis health check routine
	go db.RedisHealthCheck()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Start the Kafka consumer in a separate goroutine
	go events.ConsumeUserCreatedEvents()

	r.Get("/api/wallets/{accountId}", handlers.GetWallet)

	fmt.Println("Wallet server is running on Port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
