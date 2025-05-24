package handlers

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"github.com/AlexBetita/go_prac/internal/middlewares"
	"github.com/AlexBetita/go_prac/internal/models"
	"github.com/AlexBetita/go_prac/internal/services"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FolderHandler struct {
	service *services.FolderService
}

func NewFolderHandler(s *services.FolderService) *FolderHandler {
	return &FolderHandler{service: s}
}

func (h *FolderHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middlewares.User(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var folder models.Folder
	if err := json.NewDecoder(r.Body).Decode(&folder); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	folder.CreatedBy = user.ID

	if err := h.service.CreateFolder(r.Context(), &folder); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(folder)
}

func (h *FolderHandler) GetById(w http.ResponseWriter, r *http.Request) {
	folderIDParam := chi.URLParam(r, "id")
	folderID, err := primitive.ObjectIDFromHex(folderIDParam)
	if err != nil {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	folder, err := h.service.GetFolder(r.Context(), folderID)
	if err != nil {
		http.Error(w, "Folder not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(folder)
}

func (h *FolderHandler) Get(w http.ResponseWriter, r *http.Request) {
	user := middlewares.User(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	skip := (page - 1) * limit

	folders, total, err := h.service.GetFoldersPaginated(r.Context(), user.ID, limit, skip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"folders":     folders,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": int(math.Ceil(float64(total) / float64(limit))),
	})
}

func (h *FolderHandler) Update(w http.ResponseWriter, r *http.Request) {
	user := middlewares.User(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		ID      string         `json:"id"`
		Update  map[string]any `json:"update"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	oid, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateFolder(r.Context(), oid, req.Update); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}


func (h *FolderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user := middlewares.User(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	oid, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteFolder(r.Context(), oid); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}


func (h *FolderHandler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	user := middlewares.User(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	folderID, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	// delete all interactions that belong to this folder
	if err := h.service.DeleteInteractionsByFolder(r.Context(), folderID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// delete folder
	if err := h.service.DeleteFolder(r.Context(), folderID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}


	w.WriteHeader(http.StatusNoContent)
}

func (h *FolderHandler) ToggleFavorite(w http.ResponseWriter, r *http.Request) {
	user := middlewares.User(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		ID       string `json:"id"`
		Favorite bool   `json:"favorite"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	oid, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	if err := h.service.ToggleFavorite(r.Context(), oid, req.Favorite); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *FolderHandler) GetFavorites(w http.ResponseWriter, r *http.Request) {
	user := middlewares.User(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	skip := (page - 1) * limit

	folders, total, err := h.service.GetFavoriteFoldersPaginated(r.Context(), user.ID, limit, skip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"folders":     folders,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": int(math.Ceil(float64(total) / float64(limit))),
	})
}


