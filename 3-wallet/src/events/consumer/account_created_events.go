package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"example.com/m/v2/src/db"
	"example.com/m/v2/src/events"
	"example.com/m/v2/src/models"
	eventsModel "example.com/m/v2/src/models/events"
	"example.com/m/v2/src/utils"
	"github.com/segmentio/kafka-go"
)

//TODO: Add a Dead Message Queue

// func createReaderWithRetries(ctx context.Context) (*kafka.Reader, error) {

// 	maxRetries := 5
// 	backoff := time.Second * 2

// 	for i := 0; i < maxRetries; i++ {
// 		log.Println("Hit")
// 		reader := kafka.NewReader(kafka.ReaderConfig{
// 			Brokers:           []string{"broker:9092"},
// 			GroupID:           consumerGroup,
// 			Topic:             topic,
// 			MinBytes:          10e3,
// 			MaxBytes:          10e6,
// 			MaxWait:           1 * time.Second,
// 			ReadBackoffMin:    50 * time.Millisecond,
// 			ReadBackoffMax:    2 * time.Second,
// 			HeartbeatInterval: 5 * time.Second,
// 			CommitInterval:    time.Second,
// 			StartOffset:       kafka.FirstOffset,
// 			Logger:            kafka.LoggerFunc(log.Printf),
// 			ErrorLogger:       kafka.LoggerFunc(log.Printf),
// 		})

// 		// Try to verify connection
// 		_, err := reader.FetchMessage(ctx)
// 		if err == nil || err == context.DeadlineExceeded {
// 			return reader, nil
// 		}

// 		log.Printf("Failed to connect to Kafka, attempt %d/%d: %v", i+1, maxRetries, err)
// 		reader.Close()
// 		time.Sleep(backoff)
// 		backoff *= 2 // Exponential backoff
// 	}

//		return nil, fmt.Errorf("failed to connect to Kafka after %d attempts", maxRetries)
//	}
//
// Ensure topic exists before starting consumer

func ConsumeAccountCreatedEvents(ctx context.Context) error {

	// Configure the reader with more reasonable timeouts
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:           []string{"broker:9092"},
		GroupID:           events.ConsumerGroup,
		Topic:             events.AccountCreatedTopic,
		MinBytes:          10e3,
		MaxBytes:          10e6,
		MaxWait:           3 * time.Second,
		ReadBackoffMax:    5 * time.Second,
		HeartbeatInterval: 10 * time.Second,
		SessionTimeout:    30 * time.Second,
		StartOffset:       kafka.FirstOffset,
		// Logger:                kafka.LoggerFunc(log.Printf),
		WatchPartitionChanges: true,
	})
	defer reader.Close()

	log.Printf("Started consuming topic: %s with group: %s", events.AccountCreatedTopic, events.ConsumerGroup)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				if err != context.Canceled {
					log.Printf("Error reading message: %v", err)
				}
				time.Sleep(time.Second)
				continue
			}

			if err := processMessage(msg); err != nil {
				log.Printf("Error processing message: %v", err)
				continue
			}
		}
	}
}

// func processMessageWithRetries(msg kafka.Message) error {
// 	maxRetries := 3
// 	var lastErr error

// 	for attempt := 1; attempt <= maxRetries; attempt++ {
// 		err := processMessage(msg)
// 		if err == nil {
// 			return nil
// 		}

// 		lastErr = err
// 		log.Printf("Attempt %d failed: %v", attempt, err)

// 		if attempt < maxRetries {
// 			backoff := time.Duration(attempt*attempt) * 100 * time.Millisecond
// 			time.Sleep(backoff)
// 		}
// 	}

// 	return fmt.Errorf("all attempts failed: %v", lastErr)
// }

func processMessage(msg kafka.Message) error {
	log.Printf("Processing message with key: %s", string(msg.Key))

	var event eventsModel.AccountCreatedEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %v", err)
	}

	walletId, err := utils.GenerateWalletId()
	if err != nil {
		return fmt.Errorf("failed to generate wallet ID: %v", err)
	}

	currency, err := utils.ParseCurrency(event.Currency)
	if err != nil {
		return fmt.Errorf("failed to parse currency: %v", err)
	}
	log.Printf("Creating the Wallet ID#%s", walletId)
	log.Printf("Creating the Account ACC#%s", event.AccountId)

	var walletStatus models.WalletStatus
	if event.KYCStatus == "pending" {
		walletStatus = models.WalletStatusInactive
	} else {
		walletStatus = models.WalletStatusActive
	}

	wallet := &models.Wallet{
		WalletId:     walletId,
		AccountId:    event.AccountId,
		Type:         models.WalletTypePrimary,
		Balance:      0,
		Currency:     currency,
		Status:       walletStatus,
		IsDefault:    true,
		DailyLimit:   5000,
		MonthlyLimit: 20000,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := db.AddWallet(wallet); err != nil {
		return fmt.Errorf("failed to create wallet: %v", err)
	}

	if err := db.SetDefaultWallet(event.AccountId, wallet.WalletId); err != nil {
		return fmt.Errorf("failed to set default wallet: %v", err)
	}

	log.Printf("Successfully processed wallet creation for account: %s", event.AccountId)
	return nil
}
