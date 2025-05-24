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
	"github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRouter(cfg *config.Config, mongoClient *mongo.Client, oaClient *openai.Client,) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.Recoverer)

	db       := mongoClient.Database(cfg.DBName)
	userRepo := repositories.NewUserRepository(db)
	postRepo := repositories.NewPostRepository(db)
	interactionRepo := repositories.NewInteractionRepository(db)
	folderRepo := repositories.NewFolderRepository(db)
	messageRepo := repositories.NewMessageRepository(db)

	authSvc  := services.NewAuthService(userRepo, cfg.JWTSecret)
	botSvc   := services.NewBotService(postRepo, interactionRepo, oaClient)
	postSvc := services.NewPostService(postRepo, oaClient)
	interactionSvc := services.NewInteractionService(interactionRepo, messageRepo)
	folderSvc := services.NewFolderService(folderRepo, interactionRepo, messageRepo)
	messageSvc := services.NewMessageService(messageRepo)

	autH  := handlers.NewAuthHandler(authSvc)
	botH := handlers.NewBotHandler(botSvc)
	postH := handlers.NewPostHandler(postSvc)
	interactionH := handlers.NewInteractionHandler(interactionSvc)
	folderH := handlers.NewFolderHandler(folderSvc)
	messageH := handlers.NewMessageHandler(messageSvc)

	authMW := middlewares.Auth(cfg.JWTSecret, authSvc)

	r.Route("/api", func(api chi.Router) {
		MountAuthRoutes(api, cfg, autH, userRepo, authMW)
		MountBotRoutes(api, botH, authMW)
		MountPostRoutes(api, postH)
		MountInteractionRoutes(api, interactionH, authMW)
		MountFolderRoutes(api, folderH, authMW)
		MountMessageRoutes(api, messageH, authMW)
	})

	return r
}