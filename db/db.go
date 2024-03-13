package db

import (
	"context"
	"fmt"

	// "fmt"
	"log"
	// "time"

	"github.com/JusAeng/manga-tracker-api-go/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
  
func Connect() {
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	// defer cancel()

	// client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://planc:VQjaiVg24ZUeVTR7@cluster0.ntxuynn.mongodb.net/?retryWrites=true&w=majority"))
	DB_User,err := config.GetEnv("DB_USER")
	if err != nil{
		log.Fatalf("Fail to load env")
	}
	DB_Pass,err := config.GetEnv("DB_PASS")
	if err != nil{
		log.Fatalf("Fail to load env")
	}
	dbConnectionString := fmt.Sprintf("mongodb+srv://%s:%s@cluster0.ntxuynn.mongodb.net/?retryWrites=true&w=majority", DB_User, DB_Pass)
	clientOptions := options.Client().ApplyURI(dbConnectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to cluster: %v", err)
	}

	// err = client.Connect(context.TODO())
	// if err != nil {
	// 	log.Fatalf("Failed to connect to cluster: %v", err)
	// }

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Failed to ping cluster: %v", err)
	}

	Client = client
	log.Printf("Connected to MongoDB!")
}
 