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
		// All routes need authentication, so place them in the authenticated group
		folder.Group(func(prFolder chi.Router) {
			prFolder.Use(authMW)

			// Folder management routes
			prFolder.Get("/", fh.Get)                // Get all folders
			prFolder.Get("/:id", fh.GetById)         // Get folder by ID
			prFolder.Post("/", fh.Create)            // Create a new folder
			prFolder.Put("/", fh.Update)             // Update an existing folder
			prFolder.Delete("/", fh.Delete)          // Delete a specific folder
			prFolder.Delete("/deleteAll", fh.DeleteAll) // Delete all folders and interactions

			// Favorites management routes
			prFolder.Get("/favorites", fh.GetFavorites) // Get user favorites
			prFolder.Post("/favorite", fh.ToggleFavorite) // Toggle favorite folder status
		})
	})
}