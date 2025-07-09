package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
)

type GitHubEvent struct {
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
	Pusher struct {
		Name string `json:"name"`
	} `json:"pusher"`
}

func main() {
	// PostgreSQL setup
	dbURL := "postgres://devsync:devsyncpass@localhost:5432/devsyncdb?sslmode=disable"
	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to Postgres: %v", err)
	}
	defer dbpool.Close()
	fmt.Println("✅ Connected to Postgres")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     "devsync.events.raw",
		GroupID:   "event-processor-group",
		MinBytes:  1,
		MaxBytes:  10e6,
		MaxWait:   1 * time.Second,
		StartOffset: kafka.LastOffset, // Only process new messages
	})

	fmt.Println("🔁 Listening for events on Kafka topic: devsync.events.raw")

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("❌ Error reading message: %v", err)
		}

		var event GitHubEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("⚠️ Invalid JSON: %v\n", err)
			continue
		}

		// Save to Postgres
		_, err = dbpool.Exec(context.Background(), `
			INSERT INTO github_events (repo_name, pusher_name)
			VALUES ($1, $2)
		`, event.Repository.FullName, event.Pusher.Name)

		if err != nil {
			log.Printf("❌ Failed to insert into DB: %v\n", err)
		} else {
			fmt.Printf("✅ Saved → Repo: %s, Pusher: %s\n", event.Repository.FullName, event.Pusher.Name)
		}
	}
}
