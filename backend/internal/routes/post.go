package routes

import (
	"net/http"

	"github.com/AlexBetita/go_prac/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func MountPostRoutes(r chi.Router, 
	postH *handlers.PostHandler, 
	authMW func(http.Handler) http.Handler,
) {
	r.With(authMW).Get("/posts/{id}", postH.GetBlogByID)
}