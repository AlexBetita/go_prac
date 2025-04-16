package main

import (
	"context"
	"log"
	"os"

	"github.com/AlexBetita/go_prac/internal/config"
	"github.com/AlexBetita/go_prac/internal/db"
	"github.com/AlexBetita/go_prac/internal/db/seed"
	openai "github.com/sashabaranov/go-openai"
)

func main() {
	ctx := context.Background()
	cfg := config.New()

	client := db.Connect(cfg)
	defer client.Disconnect(ctx)

	db := client.Database(cfg.DBName)
	openaiClient := openai.NewClient(cfg.OpenAIKey)

	switch os.Getenv("SEED_MODE") {
	case "embed":
		log.Println("🔁 Running EmbedSeededPosts...")
		if err := seed.EmbedSeededPosts(ctx, db, openaiClient); err != nil {
			log.Fatal("❌ Embed failed:", err)
		}
	default:
		log.Println("🌱 Running SeedPosts...")
		if err := seed.SeedPosts(ctx, db); err != nil {
			log.Fatal("❌ Seed failed:", err)
		}
	}
}
