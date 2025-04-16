package routes

import (
	"net/http"

	"github.com/AlexBetita/go_prac/internal/config"
	"github.com/AlexBetita/go_prac/internal/handlers"
	"github.com/AlexBetita/go_prac/internal/oauth"
	"github.com/AlexBetita/go_prac/internal/repositories"
	"github.com/go-chi/chi/v5"
)

func MountAuthRoutes(
	r chi.Router,
	cfg *config.Config,
	authH *handlers.AuthHandler,
	userRepo repositories.UserRepository,
	authMW func(http.Handler) http.Handler,
) {

	r.Route("/auth", func(auth chi.Router) {
		auth.Post("/register", authH.Register)
		auth.Post("/login",    authH.Login)

		auth.Group(func(pr chi.Router) {
			pr.Use(authMW)
			pr.Get("/profile", authH.Profile)
		})

		if cfg.GoogleClientID != "" {
			googleH := oauth.NewGoogleHandler(cfg, userRepo)
			auth.Get("/google/login",    googleH.Login)
			auth.Get("/google/callback", googleH.Callback)
		}
	})
}
