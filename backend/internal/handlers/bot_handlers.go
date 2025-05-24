package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/AlexBetita/go_prac/internal/bot"
	"github.com/AlexBetita/go_prac/internal/middlewares"
	"github.com/AlexBetita/go_prac/internal/models"
	"github.com/AlexBetita/go_prac/internal/services"
	"github.com/openai/openai-go"
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

	// Call your BotService.GenerateRequest with parsed interactionID and systemPrompt
	resp, err := h.service.GenerateRequest(r.Context(), user.ID, interactionOID, fullMessage, systemPrompt)
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

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	var prompt string
	if body.SystemPrompt != nil {
		prompt = *body.SystemPrompt
	} else {
		prompt = `
		Share a brief description of what you plan to do in 40 words.
		You are pretty good at whatever you are requested to do. 
		`
	}

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(prompt),
		openai.UserMessage(body.Message),
	}

	stream := h.service.GenerateRequestStream(r.Context(), user.ID, interactionOID, messages, 
	[]string{"get_related_posts", 
	"create_blog_post"})

	acc := openai.ChatCompletionAccumulator{}

	for stream.Next() {
		chunk := stream.Current()
		acc.AddChunk(chunk)

		// When this fires, the current chunk value will not contain content data
		if _, ok := acc.JustFinishedContent(); ok {
			println()
			println("finish-event: Content stream finished")
		}

		if refusal, ok := acc.JustFinishedRefusal(); ok {
			println()
			println("finish-event: refusal stream finished:", refusal)
			println()
		}

		if tool, ok := acc.JustFinishedToolCall(); ok {
			var toolStrOutput string
			// Run the tool handler function to get output bytes
			ctxWithValues := h.enrichContext(r.Context(), user.ID)

			// Check if tool needs a follow up
			if tool.Name == "get_related_posts" {
				spec, found := bot.Registry[tool.Name]
				if !found {
					http.Error(w, "unknown tool", http.StatusInternalServerError)
					return
				}
				
				// Parse arguments JSON into raw message
				args := json.RawMessage(tool.Arguments)
				
				out, err := spec.Handle(ctxWithValues, args)

				if err != nil {
					http.Error(w, fmt.Sprintf("tool handler error: %v", err), http.StatusInternalServerError)
					return
				}

				payloadBytes, ok := out.([]byte)
				if !ok {
					http.Error(w, "invalid tool output type", http.StatusInternalServerError)
					return
				}

				payloadStr := string(payloadBytes)
				relatedPostsSystemMsg := openai.SystemMessage(`You just fetched related blog post data.

				Please format it as follows:

				1. Keep only these fields: slug, title, views, created_by.
				2. Start with a short, friendly intro (e.g. “Hey there! I found some posts you might like…”).
				3. Use **clean, well-formatted Markdown tables**. One table for relevant posts, another for not relevant.
				4. Use meaningful section headings like "### 🚀 Programming Posts" and "### 🌴 Other Interesting Reads".
				5. Keep the tone conversational but concise. Add a few light emojis for charm, but don’t overdo it.
				6. Keep spacing and formatting neat for maximum clarity.
				7. No explanations or bullet lists — only the intro and two tables.

				Thanks!`)

				followupMessages := []openai.ChatCompletionMessageParamUnion{
					acc.Choices[0].Message.ToParam(),
					openai.ToolMessage(payloadStr, acc.Choices[0].Message.ToolCalls[0].ID),
					relatedPostsSystemMsg,
					openai.UserMessage(body.Message),
				}

				followupStream := h.service.GenerateRequestStream(
					r.Context(), user.ID, interactionOID,
					followupMessages,
					nil,
				)

				followupAcc := openai.ChatCompletionAccumulator{}

				for followupStream.Next() {
					followupChunk := followupStream.Current()
					followupAcc.AddChunk(followupChunk)
					if len(followupChunk.Choices) > 0 {
						text := followupChunk.Choices[0].Delta.Content
						if len(text) > 0 {
							fmt.Fprintf(w, "event: followup_chunk\ndata: %s\n\n", escapeSSE(text))
							flusher.Flush()
						}
					}
				}

				if err := followupStream.Err(); err != nil {
					fmt.Fprintf(w, "event: error\ndata: %s\n\n", escapeSSE(err.Error()))
					flusher.Flush()
					return
				}
				finalContent, _ := followupAcc.JustFinishedContent()
				toolStrOutput = string(finalContent)
			} else {
				spec, found := bot.Registry[tool.Name]
				if !found {
					http.Error(w, "unknown tool", http.StatusInternalServerError)
					return
				}
				
				// Parse arguments JSON into raw message
				args := json.RawMessage(tool.Arguments)
				out, err := spec.Handle(ctxWithValues, args)
				if err != nil {
					http.Error(w, fmt.Sprintf("tool handler error: %v", err), http.StatusInternalServerError)
					return
				}

				switch v := out.(type) {
				case []byte:
					toolStrOutput = string(v)
				case string:
					toolStrOutput = v
				default:
					// Try to marshal to JSON if it's a struct or something else
					b, err := json.Marshal(v)
					if err != nil {
						http.Error(w, "failed to marshal tool output", http.StatusInternalServerError)
						return
					}
					toolStrOutput = string(b)
				}

			}

			payloadMap := map[string]interface{}{
				"index": tool.Index,
				"name": tool.Name,
				"arguments": json.RawMessage(tool.Arguments),
				"content": toolStrOutput,
			}
			payloadBytes, _ := json.Marshal(payloadMap)
			if _, err := fmt.Fprintf(w, "event: tool_call\ndata: %s\n\n", escapeSSE(string(payloadBytes))); err != nil {
				log.Println("write error:", err)
				return
			}
			flusher.Flush()
		}

		if len(chunk.Choices) > 0 {
			text := chunk.Choices[0].Delta.Content
			if len(text) > 0 {
				fmt.Fprintf(w, "event: chunk\ndata: %s\n\n", escapeSSE(text))
				flusher.Flush()
			}
		}
		
	}

	if err := stream.Err(); err != nil {
		fmt.Fprintf(w, "event: error\ndata: %s\n\n", escapeSSE(err.Error()))
		flusher.Flush()
		return
	}

	if acc.Usage.TotalTokens > 0 {
		tokenInt := acc.Usage.TotalTokens
		tokenStr := strconv.FormatInt(tokenInt, 10)
		fmt.Fprintf(w, "event: metadata\ndata: %s\n\n", escapeSSE(tokenStr))
	}

	// Send done event with final content
	finalContent, _ := acc.JustFinishedContent()
	finalJSON, _ := json.Marshal(map[string]string{"final": finalContent})
	fmt.Fprintf(w, "event: done\ndata: %s\n\n", finalJSON)
	flusher.Flush()
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