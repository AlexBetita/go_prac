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

type MessageHandler struct {
	service *services.MessageService
}

func NewMessageHandler(s *services.MessageService) *MessageHandler {
	return &MessageHandler{service: s}
}

func (h *MessageHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middlewares.User(r.Context())
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var msg models.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	msg.UserID = user.ID

	if err := h.service.Create(r.Context(), &msg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(msg)
}

func (h *MessageHandler) GetByInteraction(w http.ResponseWriter, r *http.Request) {
	interactionIDParam := chi.URLParam(r, "interactionId")
	interactionID, err := primitive.ObjectIDFromHex(interactionIDParam)
	if err != nil {
		http.Error(w, "Invalid interaction ID", http.StatusBadRequest)
		return
	}

	messages, err := h.service.GetByInteraction(r.Context(), interactionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
