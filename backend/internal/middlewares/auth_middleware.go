package middlewares

import (
	"context"
	"net/http"
	"strings"
	
	"github.com/AlexBetita/go_prac/internal/models"
	"github.com/AlexBetita/go_prac/pkg/utils"
	"github.com/AlexBetita/go_prac/internal/services"
)

type ctxKey string

const userKey ctxKey = "user"

func User(ctx context.Context) *models.User {
	u, _ := ctx.Value(userKey).(*models.User)
	return u
}

func Auth(jwtSecret string, authSvc services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "missing or malformed Authorization header", http.StatusUnauthorized)
				return
			}

			userID, err := utils.ValidateJWT(parts[1], jwtSecret)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			user, err := authSvc.GetUserByID(r.Context(), userID)
			if err != nil {
				http.Error(w, "user not found", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}