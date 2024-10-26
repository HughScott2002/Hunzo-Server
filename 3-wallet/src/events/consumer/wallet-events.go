package consumer

import (
	"context"
	"encoding/json"
	"log"

	"example.com/m/v2/src/models/events"
	"github.com/segmentio/kafka-go"
)

func ConsumeAccountCreatedEvents() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"broker:9092"},
		Topic:     "account-created",
		GroupID:   "wallet-service",
		Partition: 0,
	})

	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		var event events.AccountCreatedEvent
		err = json.Unmarshal(msg.Value, &event)
		// AccountId, KYCStatus
		if err != nil {
			log.Printf("Error unmarshaling event: %v", err)
			continue
		}

		// db.AddWallet()
	}
}
