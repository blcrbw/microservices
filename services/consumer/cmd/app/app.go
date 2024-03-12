package main

import (
	"consumer/internal/config"
	"consumer/internal/pkg/client/postgresql"
	"consumer/internal/test_records"
	repo "consumer/internal/test_records/db"
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	testTopic = "test"
	//userTopic    = "user"
	//productTopic = "product"

	broker0Address = "kafka-0:9092"
	broker1Address = "kafka-1:9092"
	broker2Address = "kafka-2:9092"
)

func main() {
	log.Println("Consumer started!")
	ctx, cancel := context.WithCancel(context.Background())
	pgPoolConfig := config.GetDbConfig()

	numCPU := 8

	mesChan := make(chan test_records.TestRecord, numCPU)

	db, err := postgresql.NewClient(ctx, 5, *pgPoolConfig)
	if err != nil {
		cancel()
		log.Fatalln("Unable to create connection pool:", err)
	}
	log.Println("DB pool connection created!")
	rep := repo.NewRepository(db)

	wgWriter := &sync.WaitGroup{}
	for i := 0; i <= numCPU/2; i++ {
		wgWriter.Add(1)
		go func() {
			defer wgWriter.Done()
			psqlWriter(ctx, mesChan, rep, i)
		}()
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker0Address, broker1Address, broker2Address},
		GroupID:  "consumer-test-group-id",
		Topic:    testTopic,
		MaxBytes: 10e6,
	})
	log.Println("Kafka reader started!")

	wgReader := &sync.WaitGroup{}
	for i := 0; i <= numCPU/2; i++ {
		wgReader.Add(1)
		go func() {
			defer wgReader.Done()
			kafkaReader(ctx, mesChan, r, i)
		}()
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigChan:
		if err := r.Close(); err != nil {
			log.Println("failed to close reader:", err)
		}
		cancel()

		return
	}
}

func psqlWriter(ctx context.Context, mesChan <-chan test_records.TestRecord, rep test_records.Repository, id int) {
	records := make([]*test_records.TestRecord, 0)
	defer writeMultiple(ctx, rep, records, id)
	for {
		select {
		case <-ctx.Done():
			return
		case rec, ok := <-mesChan:
			if !ok {
				return
			}
			records = append(records, &rec)
		case <-time.After(time.Millisecond * 50):
			if len(records) > 0 {
				writeMultiple(ctx, rep, records, id)
				records = make([]*test_records.TestRecord, 0)
			}
		}
	}
}

func writeMultiple(ctx context.Context, rep test_records.Repository, records []*test_records.TestRecord, id int) {
	if len(records) > 0 {
		if err := rep.CreateMultiple(ctx, records); err == nil {
			log.Printf("Thread %d: %d Test records written: %s\n", id, len(records), records)
		} else {
			log.Printf("Thread %d: Internal server/psql error: %s\n", id, err)
		}
	} else {
		log.Printf("Thread %d: Nothing to write.\n", id)
	}
}

func kafkaReader(ctx context.Context, mesChan chan<- test_records.TestRecord, r *kafka.Reader, id int) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			m, err := r.ReadMessage(ctx)
			if err != nil {
				break
			}

			now := time.Now()
			rec := test_records.TestRecord{
				Text:    string(m.Value),
				Created: now,
				Stored:  now,
			}
			mesChan <- rec
			log.Printf("Thread %d: Added to the channel: message at topic/partition/offset %v/%v/%v: %s = %s\n", id, m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		}
	}
}
