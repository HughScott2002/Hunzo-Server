package events

import (
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func List_topics() error {
	log.Printf("Checking the Broker Connection...")
	conn, err := kafka.Dial("tcp", "broker:9092")
	if err != nil {
		panic("Error in getting connected to broker")
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		return err
	}

	m := map[string]struct{}{}

	for _, p := range partitions {
		m[p.Topic] = struct{}{}
	}
	for k := range m {
		fmt.Println(k)
	}
	return nil
}
