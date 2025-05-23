package routes

import (
	"github.com/AlexBetita/go_prac/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func MountPostRoutes(r chi.Router, 
	postH *handlers.PostHandler,
) {
	 r.Get("/posts/search", postH.SearchPosts)
	r.Get("/posts/{identifier}", postH.GetPost)
}