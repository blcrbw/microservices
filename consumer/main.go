package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
)

const (
	testTopic     = "test"
	userTopic     = "user"
	productTopic  = "product"
	
	broker0Address = "kafka-0:9092"
	broker1Address = "kafka-1:9092"
	broker2Address = "kafka-2:9092"
)


func consumeTest(ctx context.Context) {
}

func main() {
	fmt.Println("Consumer started!")
	poolConfig, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln("Unable to parse DATABASE_URL:", err)
	}

	db, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalln("Unable to create connection pool:", err)
	}
	fmt.Println("DB pool connection created!")

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{broker0Address, broker1Address, broker2Address},
		GroupID:   "consumer-test-group-id",
		Topic:     testTopic,
		MaxBytes:  10e6,
	})
	fmt.Println("Kafka reader started!")

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		if (db != nil) {
			now := time.Now()
			if _, err := db.Exec(context.Background(), `insert into test_records(text, created, stored) values ($1, $2, $3)`, string(m.Value), now, now); err == nil {
				fmt.Printf("DB record written: message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
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
