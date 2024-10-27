package main

import (
	"fmt"
	"log"
	"net/http"

	"example.com/m/v2/src/db"
	"example.com/m/v2/src/events"
	"example.com/m/v2/src/events/consumer"
	"example.com/m/v2/src/server/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	db.Init()

	runKafkaRouties := true
	err := events.List_topics()
	if err != nil {
		fmt.Printf("%s", err.Error())
		runKafkaRouties = false
	}
	if runKafkaRouties {
		// Start the Kafka consumer in a separate goroutine
		go consumer.ConsumeAccountCreatedEvents()
	}

	r.Route("/api/wallets", func(r chi.Router) {
		r.Get("/{walletId}", handlers.GetWallet)
		r.Get("/list/{accountId}", handlers.ListWallets)
	})

	fmt.Println("Wallet server is running on Port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
