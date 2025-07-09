package kafka

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

var writer *kafka.Writer

func InitProducer() {
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
			kafkaBroker = "localhost:9092"
	}

	writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{kafkaBroker},
		Topic:        "devsync.events.raw",
		Balancer:     &kafka.LeastBytes{},
		WriteTimeout: 10 * time.Second,
	})

	log.Println("✅ Kafka producer initialized")
}

func PublishMessage(value []byte) error {
	msg := kafka.Message{
			Value: value,
	}

	return writer.WriteMessages(context.Background(), msg)
}
