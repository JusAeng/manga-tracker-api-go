package db

import (
	"context"
	"log"

	"github.com/JusAeng/manga-tracker-api-go/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect() {
	databaseURL, err := config.GetEnv("DATABASE_URL")
	if err != nil || databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Failed to create connection pool: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	Pool = pool
	log.Println("Connected to PostgreSQL!")
}
