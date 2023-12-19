package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
  
func Connect() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://planc:VQjaiVg24ZUeVTR7@cluster0.ntxuynn.mongodb.net/?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to cluster: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping cluster: %v", err)
	}

	Client = client
	log.Printf("Connected to MongoDB!")
}
  