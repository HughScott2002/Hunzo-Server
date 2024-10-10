package events

import (
	"context"
	"encoding/json"
	"log"

	"example.com/m/v2/src/models/events"
	"example.com/m/v2/src/utils"
	"github.com/segmentio/kafka-go"
)

func ConsumeUserCreatedEvents() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"broker:9092"},
		Topic:     "user-created",
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

		var event events.UserCreatedEvent
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			log.Printf("Error unmarshaling event: %v", err)
			continue
		}

		utils.CreateWallet(event.AccountId)
	}
}
