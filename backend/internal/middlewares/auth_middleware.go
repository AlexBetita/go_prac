package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/AlexBetita/go_prac/pkg/utils"
)

type contextKey string

const userIDKey contextKey = "userID"

func UserIDKey() contextKey {
    return userIDKey
}

func Auth(jwtSecret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "missing auth header", http.StatusUnauthorized)
                return
            }
            parts := strings.Split(authHeader, " ")
            if len(parts) != 2 || parts[0] != "Bearer" {
                http.Error(w, "invalid auth header", http.StatusUnauthorized)
                return
            }
            userID, err := utils.ValidateJWT(parts[1], jwtSecret)
            if err != nil {
                http.Error(w, "invalid token", http.StatusUnauthorized)
                return
            }
            ctx := context.WithValue(r.Context(), userIDKey, userID)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}