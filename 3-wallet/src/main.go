package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"example.com/m/v2/src/db"
	"example.com/m/v2/src/events/consumer"
	"example.com/m/v2/src/server/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	db.Init()

	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown gracefully
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down gracefully...")
		cancel()
	}()

	// Start the consumer in a goroutine
	go func() {
		log.Println("Starting Kafka consumer...")
		if err := consumer.ConsumeAccountCreatedEvents(ctx); err != nil {
			log.Printf("Kafka consumer error: %v", err)
		}
	}()

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// // List topics before starting consumer
	// if err := events.List_topics(); err != nil {
	// 	log.Printf("Error listing topics: %v", err)
	// }

	// Start the consumer (which will now create the topic if needed)
	// go consumer.ConsumeAccountCreatedEvents(ctx)
	// go consumer.TestConsumer()
	// runKafkaRouties := true
	// err := events.List_topics()
	// if err != nil {
	// 	fmt.Printf("%s", err.Error())
	// 	runKafkaRouties = false
	// }
	// if runKafkaRouties {
	// 	ctx, cancel := context.WithCancel(context.Background())
	// 	defer cancel()

	// 	// Start the Kafka consumer in a separate goroutine
	// 	go consumer.ConsumeAccountCreatedEvents(ctx)
	// }
	// go events.List_topics()

	// go func() {

	// 	kafkaConsumer := consumer.NewKafkaConsumer([]string{"broker:9092"})
	// 	kafkaConsumer.Start()

	// 	// Set up signal handling for graceful shutdown
	// 	signals := make(chan os.Signal, 1)
	// 	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// 	// Wait for termination signal
	// 	sig := <-signals
	// 	log.Printf("Received signal: %v, initiating shutdown...", sig)

	// 	// Stop the consumer
	// 	kafkaConsumer.Stop()
	// }()

	r.Route("/api/wallets", func(r chi.Router) {
		r.Get("/{walletId}", handlers.GetWallet)
		r.Get("/list/{accountId}", handlers.ListWallets)
	})

	fmt.Println("Wallet server is running on Port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
