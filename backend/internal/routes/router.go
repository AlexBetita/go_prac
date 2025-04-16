package routes

import (
	"net/http"

	"github.com/AlexBetita/go_prac/internal/config"
	"github.com/AlexBetita/go_prac/internal/handlers"
	"github.com/AlexBetita/go_prac/internal/middlewares"
	"github.com/AlexBetita/go_prac/internal/repositories"
	"github.com/AlexBetita/go_prac/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRouter(cfg *config.Config, client *mongo.Client) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.Recoverer)

	db       := client.Database(cfg.DBName)
	userRepo := repositories.NewUserRepository(db)
	postRepo := repositories.NewPostRepository(db)

	authSvc  := services.NewAuthService(userRepo, cfg.JWTSecret)
	botSvc   := services.NewBotService(postRepo, cfg.OpenAIKey)
	postSvc := services.NewPostService(postRepo, cfg.JWTSecret)

	autH  := handlers.NewAuthHandler(authSvc)
	botH := handlers.NewBotHandler(botSvc)
	postH := handlers.NewPostHandler(postSvc)

	authMW := middlewares.Auth(cfg.JWTSecret, authSvc)

	r.Route("/api", func(api chi.Router) {
		MountAuthRoutes(api, cfg, autH, userRepo, authMW)
		MountBotRoutes(api, botH, authMW)
		MountBlogRoutes(api, postH, authMW)  
	})

	return r
}