package routes

import (
	"net/http"

	"github.com/AlexBetita/go_prac/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func MountBotRoutes(
	r chi.Router,
	botH *handlers.BotHandler, 
	authMW func(http.Handler) http.Handler,
) {
	r.Route("/bot", func(pr chi.Router) {
		pr.Use(authMW)
		pr.Post("/chat", botH.Chat)
		pr.Post("/chat/stream", botH.ChatStream)
	})
}