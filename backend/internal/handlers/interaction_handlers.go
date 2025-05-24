package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AlexBetita/go_prac/internal/middlewares"
	"github.com/AlexBetita/go_prac/internal/models"
	"github.com/AlexBetita/go_prac/internal/services"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InteractionHandler struct {
	service *services.InteractionService
}

func NewInteractionHandler(s *services.InteractionService) *InteractionHandler {
	return &InteractionHandler{service: s}
}


func (h *InteractionHandler) StartOrAppend(w http.ResponseWriter, r *http.Request) {
	user := middlewares.User(r.Context())
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		InteractionID *string           `json:"interaction_id,omitempty"`
		Message       models.Message    `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var oid *primitive.ObjectID
	if req.InteractionID != nil {
		tempID, err := primitive.ObjectIDFromHex(*req.InteractionID)
		if err != nil {
			http.Error(w, "invalid interaction ID", http.StatusBadRequest)
			return
		}
		oid = &tempID
	}

	interaction, err := h.service.StartOrAppendInteraction(r.Context(), user.ID, oid, &req.Message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(interaction)
}


func (h *InteractionHandler) Update(w http.ResponseWriter, r *http.Request) {
	user := middlewares.User(r.Context())
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	interactionIDParam := chi.URLParam(r, "id")
	interactionID, err := primitive.ObjectIDFromHex(interactionIDParam)
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateInteraction(r.Context(), interactionID, updates); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
