package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"devsync/kafka"
)

type GitHubEvent struct {
	Repository struct {
			FullName string `json:"full_name"`
	} `json:"repository"`
	Pusher struct {
			Name string `json:"name"`
	} `json:"pusher"`
}

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	var event GitHubEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
	}

	payload, err := json.Marshal(event)
	if err != nil {
		http.Error(w, "Failed to encode event", http.StatusInternalServerError)
		return
	}

	// Publish to Kafka
	if err := kafka.PublishMessage(payload); err != nil {
		http.Error(w, "Failed to send to Kafka", http.StatusInternalServerError)
		fmt.Printf("❌ Failed to send event to Kafka: %v\n", err)
		return
	}

	fmt.Printf("📤 Sent event to Kafka: %s by %s\n", event.Repository.FullName, event.Pusher.Name)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Event received"))
}
