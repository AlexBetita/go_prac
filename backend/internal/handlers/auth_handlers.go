package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AlexBetita/go_prac/internal/middlewares"
	"github.com/AlexBetita/go_prac/internal/services"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
    service services.AuthService
}

func NewAuthHandler(service services.AuthService) *AuthHandler {
    return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    user, err := h.service.Register(r.Context(), req.Email, req.Password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    tokenString, user, err := h.service.Login(r.Context(), req.Email, req.Password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
    if err != nil {
        http.Error(w, "Invalid token", http.StatusInternalServerError)
        return
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        http.Error(w, "Invalid token claims", http.StatusInternalServerError)
        return
    }

    exp, ok := claims["exp"].(float64)
    if !ok {
        http.Error(w, "Missing exp in token", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "token": tokenString,
        "exp":   int64(exp),
        "user": map[string]interface{}{
            "email": user.Email,
            "provider":  "local",
            "created_at": user.CreatedAt,
            "updated_at": user.UpdatedAt,
        },
    })
}


func (h *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) {
    user := middlewares.User(r.Context())
    if user == nil {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}