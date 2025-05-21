package routes

import (
	"net/http"

	"github.com/AlexBetita/go_prac/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func MountFolderRoutes(
	r chi.Router, 
	fh *handlers.FolderHandler,
	authMW func(http.Handler) http.Handler,
) {
		r.Route("/folders", func(folder chi.Router) {
			folder.Get("/", fh.Get)
			folder.Group(func(prFolder chi.Router) {
				prFolder.Use(authMW)
				prFolder.Post("/", fh.Create)
				prFolder.Delete("/", fh.Delete)
				prFolder.Put("/", fh.Edit)
			})
			// Add more routes here
		})
}