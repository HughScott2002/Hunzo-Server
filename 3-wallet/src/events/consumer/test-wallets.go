package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

const (
	Topic         = "account-created"
	ConsumerGroup = "wallet-service"
)

type AccountCreatedEvent struct {
	AccountId string `json:"accountId"`
	Currency  string `json:"currency"`
	KYCStatus string `json:"kycstatus"`
}

type KafkaConsumer struct {
	reader *kafka.Reader
	ctx    context.Context
	cancel context.CancelFunc
}

// NewKafkaConsumer creates a new Kafka consumer instance
func NewKafkaConsumer(brokers []string) *KafkaConsumer {
	ctx, cancel := context.WithCancel(context.Background())

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  ConsumerGroup,
		Topic:    Topic,
		MaxBytes: 10e6, // 10MB
	})

	return &KafkaConsumer{
		reader: reader,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start begins consuming messages in a separate goroutine
func (c *KafkaConsumer) Start() {
	log.Printf("Starting consumer group %s for topic %s", ConsumerGroup, Topic)

	go func() {
		defer func() {
			if err := c.reader.Close(); err != nil {
				log.Printf("Error closing reader: %v", err)
			}
		}()

		for {
			select {
			case <-c.ctx.Done():
				log.Println("Shutting down consumer...")
				return
			default:
				message, err := c.reader.ReadMessage(c.ctx)
				if err != nil {
					if err == context.Canceled {
						return
					}
					log.Printf("Error reading message: %v", err)
					continue
				}

				// Parse the message
				var event AccountCreatedEvent
				if err := json.Unmarshal(message.Value, &event); err != nil {
					log.Printf("Error unmarshaling message: %v\nMessage: %s", err, string(message.Value))
					continue
				}

				// Log the parsed event
				log.Printf("Received message at topic/partition/offset %v/%v/%v:\n"+
					"AccountID: %s\n"+
					"Currency: %s\n"+
					"KYC Status: %s",
					message.Topic, message.Partition, message.Offset,
					event.AccountId, event.Currency, event.KYCStatus)
			}
		}
	}()
}

// Stop gracefully shuts down the consumer
func (c *KafkaConsumer) Stop() {
	log.Println("Stopping Kafka consumer...")
	c.cancel()
}
