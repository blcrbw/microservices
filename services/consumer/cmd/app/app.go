package main

import (
	"consumer/internal/test_records"
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"time"

	"consumer/internal/config"
	"consumer/internal/pkg/client/postgresql"
	repo "consumer/internal/test_records/db"
)

const (
	testTopic    = "test"
	userTopic    = "user"
	productTopic = "product"

	broker0Address = "kafka-0:9092"
	broker1Address = "kafka-1:9092"
	broker2Address = "kafka-2:9092"
)

func consumeTest(ctx context.Context) {
}

func main() {
	fmt.Println("Consumer started!")
	pgPoolConfig := config.GetDbConfig()

	db, err := postgresql.NewClient(context.Background(), 5, *pgPoolConfig)
	if err != nil {
		log.Fatalln("Unable to create connection pool:", err)
	}
	fmt.Println("DB pool connection created!")
	rep := repo.NewRepository(db)

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker0Address, broker1Address, broker2Address},
		GroupID:  "consumer-test-group-id",
		Topic:    testTopic,
		MaxBytes: 10e6,
	})
	fmt.Println("Kafka reader started!")

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}

		now := time.Now()
		rec := test_records.TestRecord{
			Text:    string(m.Value),
			Created: now,
			Stored:  now,
		}
		if rep != nil {
			if err := rep.Create(context.Background(), &rec); err == nil {
				fmt.Printf("DB record written: message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
				fmt.Printf("Test record: %s\n", rec)
			} else {
				fmt.Println("Internal server error: ", err)
			}
		} else {
			fmt.Printf("DB record cannot be written: message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		}
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
