package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AlexBetita/go_prac/internal/config"
	"github.com/AlexBetita/go_prac/internal/db"
	"github.com/AlexBetita/go_prac/internal/routes"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, relying on environment variables")
    }

    cfg := config.New()

    oaClient := openai.NewClient(cfg.OpenAIKey)
    mongoClient := db.Connect(cfg)
    defer mongoClient.Disconnect(context.Background())

    routeHandler := routes.NewRouter(cfg, mongoClient, oaClient)

    srv := &http.Server{
        Addr:         ":" + cfg.ServerPort,
        Handler:      routeHandler,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    go func() {
        log.Printf("Server is running on port %s", cfg.ServerPort)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen: %s\n", err)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Server Shutdown Failed:%+v", err)
    }
    log.Println("Server exiting")
}