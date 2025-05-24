package routes

import (
	"net/http"

	"github.com/AlexBetita/go_prac/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func MountInteractionRoutes(r chi.Router, ih *handlers.InteractionHandler, authMW func(http.Handler) http.Handler) {
	r.Route("/interactions", func(ir chi.Router) {
		ir.Use(authMW)
		ir.With(authMW).Post("/start", ih.StartOrAppend)
		ir.With(authMW).Put("/{id}", ih.Update)
	})
}