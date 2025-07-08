package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
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

	// Later: Push this event to Kafka
	fmt.Printf("📥 Received event from repo: %s, by user: %s\n", event.Repository.FullName, event.Pusher.Name)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Event received"))
}
