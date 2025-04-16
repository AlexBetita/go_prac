package routes

import (
	"net/http"

	"github.com/AlexBetita/go_prac/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRouter(cfg *config.Config, client *mongo.Client) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)


    r.Route("/api", func(api chi.Router) {
        MountAuthRoutes(api, cfg, client)
		MountBotRoutes(api, cfg, client)
    })

    return r
}