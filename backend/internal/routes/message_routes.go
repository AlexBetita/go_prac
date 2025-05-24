package routes

import (
	"net/http"

	"github.com/AlexBetita/go_prac/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func MountMessageRoutes(r chi.Router, mh *handlers.MessageHandler, authMW func(http.Handler) http.Handler) {
	r.Route("/messages", func(mr chi.Router) {
		mr.Use(authMW)
		mr.Post("/", mh.Create)
		mr.Get("/interaction/{interactionId}", mh.GetByInteraction)
	})
}
