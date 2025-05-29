package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	// "log"
	"net/http"
	// "strconv"
	"strings"

	"github.com/AlexBetita/go_prac/internal/bot"
	"github.com/AlexBetita/go_prac/internal/middlewares"
	"github.com/AlexBetita/go_prac/internal/models"
	"github.com/AlexBetita/go_prac/internal/services"
	// "github.com/openai/openai-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BotHandler struct {
	service *services.BotService
	awsSvc  *services.AWSService
}

func NewBotHandler(service *services.BotService, awsSvc *services.AWSService) *BotHandler {
    return &BotHandler{
        service: service,
        awsSvc:  awsSvc,
    }
}

const (
	MaxFiles     = 5
	MaxTotalSize = 10 * 1024 * 1024 // 10 MB
)

func (h *BotHandler) Chat(w http.ResponseWriter, r *http.Request) {
	user := middlewares.User(r.Context())
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var message string
	var attachments []models.Attachment

	var interactionOID *primitive.ObjectID
	var systemPrompt *string

	contentType := r.Header.Get("Content-Type")

	switch {
	case strings.HasPrefix(contentType, "application/json"):
		var body struct {
			Message       string  `json:"message"`
			InteractionID *string `json:"interaction_id,omitempty"`
			SystemPrompt  *string `json:"system_prompt,omitempty"`
			Stream        bool    `json:"stream,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || strings.TrimSpace(body.Message) == "" {
			http.Error(w, "message required", http.StatusBadRequest)
			return
		}
		message = body.Message

		if body.InteractionID != nil {
			oid, err := primitive.ObjectIDFromHex(*body.InteractionID)
			if err != nil {
				http.Error(w, "invalid interaction_id", http.StatusBadRequest)
				return
			}
			interactionOID = &oid
		}

		systemPrompt = body.SystemPrompt

	case strings.HasPrefix(contentType, "multipart/form-data"):
		if err := r.ParseMultipartForm(MaxTotalSize + 1024); err != nil {
			http.Error(w, "bad form", http.StatusBadRequest)
			return
		}
		message = r.FormValue("message")
		if strings.TrimSpace(message) == "" {
			http.Error(w, "message required", http.StatusBadRequest)
			return
		}

		// Optional: parse interaction_id and system_prompt from form values if needed
		if val := r.FormValue("interaction_id"); val != "" {
			oid, err := primitive.ObjectIDFromHex(val)
			if err != nil {
				http.Error(w, "invalid interaction_id", http.StatusBadRequest)
				return
			}
			interactionOID = &oid
		}
		if val := r.FormValue("system_prompt"); val != "" {
			systemPrompt = &val
		}

		files := r.MultipartForm.File["files"]
		if len(files) > MaxFiles {
			http.Error(w, "too many files", http.StatusBadRequest)
			return
		}

		var totalSize int64
		for _, fh := range files {
			totalSize += fh.Size
			if fh.Size > MaxTotalSize {
				http.Error(w, "file too large", http.StatusBadRequest)
				return
			}
		}
		if totalSize > MaxTotalSize {
			http.Error(w, "total upload size too large", http.StatusBadRequest)
			return
		}

		for _, fh := range files {
			f, err := fh.Open()
			if err != nil {
				http.Error(w, "file error", http.StatusInternalServerError)
				return
			}
			defer f.Close()

			url, err := h.awsSvc.UploadFile(r.Context(), fh.Filename, f, fh.Header.Get("Content-Type"))
            if err != nil {
                http.Error(w, "upload error", http.StatusInternalServerError)
                return
            }
			isImg := strings.HasPrefix(fh.Header.Get("Content-Type"), "image/")
			attachments = append(attachments, models.Attachment{
				Name:    fh.Filename,
				URL:     url,
				Type:    fh.Header.Get("Content-Type"),
				IsImage: isImg,
			})
		}

	default:
		http.Error(w, "unsupported content type", http.StatusBadRequest)
		return
	}

	// Build full message text including attachment references
	fullMessage := message
	for _, att := range attachments {
		if att.IsImage {
			fullMessage += fmt.Sprintf("\nImage URL: %s", att.URL)
		} else {
			fullMessage += fmt.Sprintf("\n[File uploaded: %s]", att.Name)
		}
	}

	resp, err := h.service.GenerateRequest(r.Context(), user.ID, interactionOID, fullMessage, systemPrompt, []string{""})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *BotHandler) ChatStream(w http.ResponseWriter, r *http.Request) {
	user := middlewares.User(r.Context())
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var body struct {
		Message       string  `json:"message"`
		InteractionID *string `json:"interaction_id,omitempty"`
		SystemPrompt  *string `json:"system_prompt,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || strings.TrimSpace(body.Message) == "" {
		http.Error(w, "message required", http.StatusBadRequest)
		return
	}

	var interactionOID *primitive.ObjectID
	if body.InteractionID != nil {
		oid, err := primitive.ObjectIDFromHex(*body.InteractionID)
		if err != nil {
			http.Error(w, "invalid interaction_id", http.StatusBadRequest)
			return
		}
		interactionOID = &oid
	}

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")


	
	h.service.GenerateRequestStream(r.Context(), user.ID, interactionOID, body.Message,
	body.SystemPrompt,
	[]string{"get_related_posts", 
	"create_blog_post"}, w)

	
}

// escapeSSE escapes newlines per SSE spec
func escapeSSE(data string) string {
	data = strings.ReplaceAll(data, "\\", "\\\\")
	data = strings.ReplaceAll(data, "\n", "\\n")
	data = strings.ReplaceAll(data, "\r", "\\r")
	return data
}

func (h *BotHandler) enrichContext(ctx context.Context, userID primitive.ObjectID) context.Context {
    ctx = context.WithValue(ctx, bot.CtxUserID, userID)
    ctx = context.WithValue(ctx, bot.CtxRepo, h.service.PostRepo())
    ctx = context.WithValue(ctx, bot.CtxClient, h.service.OpenAIClient())
	ctx = context.WithValue(ctx, bot.CtxInput, "")
    return ctx
}