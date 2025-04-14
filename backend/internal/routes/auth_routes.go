package routes

import (
	"github.com/AlexBetita/go_prac/internal/config"
	"github.com/AlexBetita/go_prac/internal/handlers"
	"github.com/AlexBetita/go_prac/internal/middlewares"
	"github.com/AlexBetita/go_prac/internal/oauth"
	"github.com/AlexBetita/go_prac/internal/repositories"
	"github.com/AlexBetita/go_prac/internal/services"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

func MountAuthRoutes(r chi.Router, cfg *config.Config, client *mongo.Client) {
	db := client.Database(cfg.DBName)
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService)

	r.Route("/auth", func(auth chi.Router) {
		auth.Post("/register", authHandler.Register)
		auth.Post("/login", authHandler.Login)

		auth.Group(func(pr chi.Router) {
			pr.Use(middlewares.Auth(cfg.JWTSecret))
			pr.Get("/profile", authHandler.Profile)
		})

		if cfg.GoogleClientID != "" {
			googleHandler := oauth.NewGoogleHandler(cfg, userRepo)
			auth.Get("/google/login", googleHandler.Login)
			auth.Get("/google/callback", googleHandler.Callback)
		}
	})
}
