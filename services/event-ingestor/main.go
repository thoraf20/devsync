package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/gorilla/mux"
	"devsync/internal/handlers"
	"devsync/kafka"
)

func main() {
	err := godotenv.Load()
	if err != nil {
			log.Println("No .env file found, using default env values")
	}

	port := os.Getenv("PORT")
	if port == "" {
			port = "8080"
	}
	// Initialize Kafka producer
	kafka.InitProducer()

	r := mux.NewRouter()
	r.HandleFunc("/webhook", handlers.WebhookHandler).Methods("POST")

	fmt.Println("✅ Event Ingestor running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
