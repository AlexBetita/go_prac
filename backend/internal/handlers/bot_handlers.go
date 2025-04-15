package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/AlexBetita/go_prac/internal/middlewares"
	"github.com/AlexBetita/go_prac/internal/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BotHandler struct{service *services.BotService }

func NewBotHandler(service *services.BotService) *BotHandler { return &BotHandler{service} }

func (h *BotHandler) Chat(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Message string `json:"Message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || strings.TrimSpace(body.Message) == "" {
		http.Error(w, "message required", http.StatusBadRequest)
		return
	}

	userIDStr, ok := r.Context().Value(middlewares.UserIDKey()).(string)
	if !ok || userIDStr == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		http.Error(w, "invalid userID", http.StatusInternalServerError)
		return
	}

	post, err := h.service.GenerateRequest(r.Context(), userID, body.Message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}
