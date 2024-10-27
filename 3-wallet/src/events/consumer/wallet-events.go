package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"example.com/m/v2/src/db"
	"example.com/m/v2/src/models"
	"example.com/m/v2/src/models/events"
	"example.com/m/v2/src/utils"
	"github.com/segmentio/kafka-go"
)

// TODO: Send to Deal Letter Queue
// TODO: Implement Retries

func ConsumeAccountCreatedEvents() {
	log.Printf("Consuming Account Created Events")
	reader := kafka.NewReader(kafka.ReaderConfig{
		// List of Kafka brokers to connect to
		Brokers: []string{"broker:9092"},
		// Topic to consume messages from
		Topic: "account-created",
		// Consumer group ID - allows multiple consumers to coordinate consuming messages.
		// Consumers with the same GroupID share the message load
		GroupID: "wallet-service",

		// Partition: 0,
		// Minimum amount of data to fetch in each request (10KB)
		// Helps optimize network usage by batching small messages
		MinBytes: 10e3,

		// Maximum amount of data to fetch in each request (10MB)
		// Prevents overwhelming the consumer with too much data at once
		MaxBytes: 10e6,

		// How often to commit the offset of processed messages
		// Commits every second in this case
		CommitInterval: time.Second,

		// How often to check for lag between consumer and producer
		// Checks every second in this case
		ReadLagInterval: time.Second,

		// How long a consumer can be inactive before being considered dead
		// Other consumers will take over its partitions after this timeout
		SessionTimeout: 20 * time.Second,

		// How often to send heartbeats to maintain consumer group membership
		// Should be lower than SessionTimeout (typically 1/3)
		HeartbeatInterval: 6 * time.Second,

		// How long messages are retained in the topic
		// Messages older than this will be deleted
		RetentionTime: 120 * time.Hour,

		// Where to start reading messages from:
		// FirstOffset = oldest available message
		// LastOffset would be newest message
		StartOffset: kafka.FirstOffset,
	})
	// reader := kafka.NewReader(kafka.ReaderConfig{
	// 	Brokers:   []string{"broker:9092"},
	// 	Topic:     "account-created",
	// 	GroupID:   "wallet-service",
	// 	Partition: 0},
	// )

	defer reader.Close()

	for {
		ctx := context.Background()
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}
		log.Printf("event.AccountCreated")

		var event events.AccountCreatedEvent
		err = json.Unmarshal(msg.Value, &event)
		// AccountId, KYCStatus, Currency
		if err != nil {
			log.Printf("Error unmarshaling event: %v", err)
			continue
		}
		log.Printf("Genning ID")

		walletId, err := utils.GenerateWalletId()
		if err != nil {
			log.Printf("Error generating walletID")
			return
		}
		log.Printf("ParsingCurrecny")

		currency, err := utils.ParseCurrency(event.Currency)
		if err != nil {
			log.Printf("Error parsing Currency")
			return
		}
		log.Printf("Create a new wallet")
		// Create a new wallet
		wallet := &models.Wallet{
			WalletId:     walletId,
			AccountId:    event.AccountId,
			Type:         models.WalletTypePrimary,
			Balance:      0,
			Currency:     currency,
			Status:       models.WalletStatusActive,
			IsDefault:    true,
			DailyLimit:   5000, // Set appropriate default limits
			MonthlyLimit: 20000,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		log.Printf("Add wallet to db")

		// Add wallet to database
		if err := db.AddWallet(wallet); err != nil {
			fmt.Printf("failed to create wallet: %v", err)
			return
		}
		log.Printf("Set wallet as defualt")

		// Set as default wallet
		if err := db.SetDefaultWallet(event.AccountId, wallet.WalletId); err != nil {
			fmt.Printf("failed to set default wallet: %v", err)
			return
		}
	}
}

// func processAccountCreatedMessage(ctx context.Context, msg kafka.Message) bool {
// 	return false
// }

// // Retry helper function
// func retry(attempts int, fn func() error) error {
// 	var err error
// 	for i := 0; i < attempts; i++ {
// 		if err = fn(); err == nil {
// 			return nil
// 		}
// 		time.Sleep(time.Second * time.Duration(i+1))
// 	}
// 	return err
// }
