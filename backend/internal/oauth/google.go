package oauth

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/AlexBetita/go_prac/internal/config"
	"github.com/AlexBetita/go_prac/internal/models"
	"github.com/AlexBetita/go_prac/internal/repositories"
	"github.com/AlexBetita/go_prac/pkg/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleHandler struct {
    oauthConfig *oauth2.Config
    repo        repositories.UserRepository
    jwtSecret   string
}

func NewGoogleHandler(cfg *config.Config, repo repositories.UserRepository) *GoogleHandler {
    return &GoogleHandler{
        oauthConfig: &oauth2.Config{
            ClientID:     cfg.GoogleClientID,
            ClientSecret: cfg.GoogleClientSecret,
            RedirectURL:  cfg.GoogleRedirectURL,
            Scopes:       []string{"email", "profile"},
            Endpoint:     google.Endpoint,
        },
        repo:      repo,
        jwtSecret: cfg.JWTSecret,
    }
}

func (h *GoogleHandler) Login(w http.ResponseWriter, r *http.Request) {
    url := h.oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
    http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *GoogleHandler) Callback(w http.ResponseWriter, r *http.Request) {
    code := r.URL.Query().Get("code")
    token, err := h.oauthConfig.Exchange(context.Background(), code)
    if err != nil {
        http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
        return
    }
    client := h.oauthConfig.Client(context.Background(), token)
    resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
    if err != nil {
        http.Error(w, "Failed to get user info", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()
    var userInfo struct {
        ID    string `json:"id"`
        Email string `json:"email"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
        http.Error(w, "Failed to parse user info", http.StatusInternalServerError)
        return
    }

    // Find or create user
    user, err := h.repo.FindByEmail(r.Context(), userInfo.Email)
    if err != nil {
        user = &models.User{
            Email:      userInfo.Email,
            Provider:   "google",
            ProviderID: userInfo.ID,
            CreatedAt:  time.Now().Unix(),
            UpdatedAt:  time.Now().Unix(),
        }
        if err := h.repo.Create(r.Context(), user); err != nil {
            http.Error(w, "Failed to create user", http.StatusInternalServerError)
            return
        }
    }

    jwtToken, _ := utils.GenerateJWT(user.ID.Hex(), h.jwtSecret)
    http.SetCookie(w, &http.Cookie{
        Name:     "token",
        Value:    jwtToken,
        Expires:  time.Now().Add(24 * time.Hour),
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteLaxMode,
    })
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"token": jwtToken})
}