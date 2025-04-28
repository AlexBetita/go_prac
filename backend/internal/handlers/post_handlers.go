package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/AlexBetita/go_prac/internal/models"
	"github.com/AlexBetita/go_prac/internal/services"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostHandler struct {
    service services.PostService
}

func NewPostHandler(service services.PostService) *PostHandler {
    return &PostHandler{service: service}
}


func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	identifier := chi.URLParam(r, "identifier")
	if identifier == "" {
		http.Error(w, "missing identifier", http.StatusBadRequest)
		return
	}

	post, err := h.service.GetPost(r.Context(), identifier)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, "post not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func (h *PostHandler) SearchPosts(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    q := r.URL.Query().Get("q")
    if q == "" {
        http.Error(w, "`q` query‑param is required", http.StatusBadRequest)
        return
    }

    limit := int64(10)
    if l := r.URL.Query().Get("limit"); l != "" {
        if v, err := strconv.ParseInt(l, 10, 64); err == nil {
            limit = v
        }
    }

    useVector := r.URL.Query().Get("vector") == "true"

    var posts []*models.Post
    var err error
    if useVector {
        posts, err = h.service.SearchPostsByVector(ctx, q, limit)
    } else {
        posts, err = h.service.SearchPosts(ctx, q, limit)
    }
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(posts)
}