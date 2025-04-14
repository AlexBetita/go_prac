package db

import (
	"context"
	"log"
	"time"

	"github.com/AlexBetita/go_prac/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(cfg *config.Config) *mongo.Client {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }
    return client
}