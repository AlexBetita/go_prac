package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AlexBetita/go_prac/internal/models"
	"github.com/AlexBetita/go_prac/internal/services"
	"github.com/AlexBetita/go_prac/internal/middlewares"
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

func (h *FolderHandler) Get(w http.ResponseWriter, r *http.Request) {
	//
}

func (h *FolderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	//
}

func (h *FolderHandler) Edit(w http.ResponseWriter, r *http.Request) {
	//
}