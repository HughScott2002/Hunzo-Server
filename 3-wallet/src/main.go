package main

import (
	"fmt"
	"log"
	"net/http"

	"example.com/m/v2/src/events/consumer"
	"example.com/m/v2/src/server/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Start the Kafka consumer in a separate goroutine
	go consumer.ConsumeAccountCreatedEvents()

	r.Route("/api/wallets", func(r chi.Router) {
		r.Get("/{accountId}", handlers.GetWallet)
	})

	fmt.Println("Wallet server is running on Port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
