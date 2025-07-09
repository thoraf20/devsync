package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
)

type RawGitHubEvent  struct {
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
	Pusher struct {
		Name string `json:"name"`
	} `json:"pusher"`
}

type GitHubEvent struct {
	ID         int       `json:"id"`
	RepoName   string    `json:"repo_name"`
	PusherName string    `json:"pusher_name"`
	ReceivedAt time.Time `json:"received_at"`
}

var dbpool *pgxpool.Pool

func handleGetEvents(w http.ResponseWriter, r *http.Request) {
	repo := r.URL.Query().Get("repo")
	pusher := r.URL.Query().Get("pusher")

	// Handle pagination defaults
	limit := 50
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	query := `
		SELECT id, repo_name, pusher_name, received_at
		FROM github_events
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	if repo != "" {
		query += fmt.Sprintf(" AND repo_name = $%d", argIndex)
		args = append(args, repo)
		argIndex++
	}

	if pusher != "" {
		query += fmt.Sprintf(" AND pusher_name = $%d", argIndex)
		args = append(args, pusher)
		argIndex++
	}

	query += fmt.Sprintf(" ORDER BY received_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := dbpool.Query(context.Background(), query, args...)
	if err != nil {
		http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	var events []GitHubEvent
	for rows.Next() {
		var e GitHubEvent
		if err := rows.Scan(&e.ID, &e.RepoName, &e.PusherName, &e.ReceivedAt); err != nil {
			log.Println("Scan error:", err)
			continue
		}
		events = append(events, e)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}


	func main() {
	dbURL := "postgres://devsync:devsyncpass@localhost:5432/devsyncdb?sslmode=disable"
	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to Postgres: %v", err)
	}
	defer dbpool.Close()
	fmt.Println("✅ Connected to Postgres")

	// Kafka consumer
	go func() {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:   []string{"localhost:9092"},
			Topic:     "devsync.events.raw",
			GroupID:   "event-processor-group",
			MinBytes:  1,
			MaxBytes:  10e6,
			MaxWait:   1 * time.Second,
			// StartOffset: kafka.LastOffset, // Only process new messages
		})
		defer reader.Close()

		fmt.Println("🔁 Listening for events on Kafka topic: devsync.events.raw")

		for {
			msg, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("❌ Error reading message: %v", err)
				continue
			}

			var event GitHubEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Printf("⚠️ Invalid JSON: %v\n", err)
				continue
			}

			_, err = dbpool.Exec(context.Background(), `
				INSERT INTO github_events (repo_name, pusher_name)
				VALUES ($1, $2)
			`, event.RepoName, event.PusherName)

			if err != nil {
				log.Printf("❌ Failed to insert into DB: %v\n", err)
			} else {
				fmt.Printf("✅ Saved → Repo: %s, Pusher: %s\n", event.RepoName, event.PusherName)
			}
		}
	}()

	// HTTP server
	http.HandleFunc("/events", handleGetEvents)
	fmt.Println("🌐 REST API running at http://localhost:8081/events")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
