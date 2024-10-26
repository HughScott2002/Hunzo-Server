package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"example.com/m/v2/src/models/events"
	"github.com/segmentio/kafka-go"
)

const (
	userCreatedTopic = "account-created"
	kafkaBroker      = "broker:9092"
	kafkaPartition   = 0
)

func ProduceUserCreatedEvent(event events.AccountCreatedEvent) error {
	conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaBroker, userCreatedTopic, kafkaPartition)
	if err != nil {
		return fmt.Errorf("failed to dail leader: %v", err)
	}
	// defer log.Println(conn.Close())
	defer conn.Close()

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %v", err)
	}
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: eventJSON},
	)
	if err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}

	return nil

}
