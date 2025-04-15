package routes

import (
	"github.com/AlexBetita/go_prac/internal/config"
	"github.com/AlexBetita/go_prac/internal/handlers"
	"github.com/AlexBetita/go_prac/internal/middlewares"
	"github.com/AlexBetita/go_prac/internal/repositories"
	"github.com/AlexBetita/go_prac/internal/services"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

func MountBotRoutes(r chi.Router, cfg *config.Config, client *mongo.Client) {
	db := client.Database(cfg.DBName)

	postRepo := repositories.NewPostRepository(db)
	botSvc   := services.NewBotService(postRepo, cfg.OpenAIKey)
	botHdl   := handlers.NewBotHandler(botSvc)

	r.Route("/bot", func(pr chi.Router) {
		pr.Use(middlewares.Auth(cfg.JWTSecret))
		pr.Post("/chat", botHdl.Chat)
	})
}