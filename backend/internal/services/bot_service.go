package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/AlexBetita/go_prac/internal/bot"
	"github.com/AlexBetita/go_prac/internal/models"
	"github.com/AlexBetita/go_prac/internal/repositories"
	openai "github.com/openai/openai-go"

	// "github.com/openai/openai-go/packages/ssestream"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BotService struct {
	pRepo   repositories.PostRepository
	iaRepo  repositories.InteractionRepository
	iaSvc   *InteractionService
	oaClient *openai.Client
}

func NewBotService(pRepo repositories.PostRepository,
	iaRepo repositories.InteractionRepository,
	msgRepo repositories.MessageRepository,
	oaClient *openai.Client) *BotService {
	iaSvc := NewInteractionService(iaRepo, msgRepo)
	return &BotService{pRepo: pRepo, iaRepo: iaRepo,
        iaSvc: iaSvc, oaClient: oaClient}
}

func (s *BotService) PostRepo() repositories.PostRepository {
    return s.pRepo
}

func (s *BotService) OpenAIClient() *openai.Client {
    return s.oaClient
}

func (s *BotService) GenerateRequest(
	ctx context.Context,
	userID primitive.ObjectID,
	interactionID *primitive.ObjectID,
	message string,
	systemPrompt *string,
	plugins []string,
) (*models.BotResponse, error) {
	
	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = openai.ChatModelGPT4o
	}

	var prompt string
	if systemPrompt != nil {
		prompt = *systemPrompt
	} else {
		prompt = "You are pretty good at whatever you are requested to do."
	}

	ctx = context.WithValue(ctx, bot.CtxUserID, userID)
	ctx = context.WithValue(ctx, bot.CtxInput, message)
	ctx = context.WithValue(ctx, bot.CtxRepo, s.pRepo)
	ctx = context.WithValue(ctx, bot.CtxClient, s.oaClient)

	// Build tools param from registry
	tools := make([]openai.ChatCompletionToolParam, 0, len(bot.Registry))
	for _, spec := range bot.Registry {
		tools = append(tools, openai.ChatCompletionToolParam{
			Function: openai.FunctionDefinitionParam{
				Name:        spec.Definition.Name,
				Description: openai.String(spec.Definition.Description),
				Parameters:  spec.Definition.Parameters,
			},
		})
	}

	// Initial messages slice with system + user
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(prompt),
		openai.UserMessage(message),
	}

	// Prepare first chat request params
	params := openai.ChatCompletionNewParams{
		Model:  model,
		Messages: messages,
		Tools:  tools,
		Seed:   openai.Int(0),
	}

	// Send first request
	firstResp, err := s.oaClient.Chat.Completions.New(ctx, params)
	if err != nil {
		return nil, err
	}
	if len(firstResp.Choices) == 0 {
		return nil, errors.New("no choices returned by OpenAI")
	}

	choiceMsg := firstResp.Choices[0].Message
	toolCalls := choiceMsg.ToolCalls

	// No function calls? Return bot response as interaction
	if len(toolCalls) == 0 {
		messageModel := &models.Message{
			UserID:      userID,
			UserContent: message,
			AssistantContent: choiceMsg.Content,
		}
		if _, err := s.iaSvc.StartOrAppendInteraction(ctx, 
			userID, 
			interactionID,
			messageModel,
			model,
			&prompt,
			); err != nil {
			return nil, err
		}
		return &models.BotResponse{
			Type:     "interaction",
			Response: choiceMsg.Content,
		}, nil
	}

	// There is a tool call — handle first tool call only for now
	call := toolCalls[0]
	spec, ok := bot.Registry[call.Function.Name]
	if !ok {
		return nil, fmt.Errorf("unknown function: %s", call.Function.Name)
	}

	// Parse tool call arguments
	var args json.RawMessage = []byte(call.Function.Arguments)
	out, err := spec.Handle(ctx, args)
	if err != nil {
		return nil, err
	}
	
	switch call.Function.Name {
	case "get_related_posts":
		payloadBytes, ok := out.([]byte)
		if !ok {
			return nil, errors.New("expected []byte output from get_related_posts handler")
		}
		payloadStr := string(payloadBytes)

		// System message guiding formatting
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

		// Append function response message
		baseMessages := []openai.ChatCompletionMessageParamUnion{
			choiceMsg.ToParam(),
			openai.ToolMessage(payloadStr, call.ID),
			relatedPostsSystemMsg,
			openai.UserMessage(message),
		}

		// Second request with function result
		followupParams := openai.ChatCompletionNewParams{
			Model:    model,
			Messages: baseMessages,
		}

		finalResp, err := s.oaClient.Chat.Completions.New(ctx, followupParams)
		if err != nil {
			return nil, err
		}
		if len(finalResp.Choices) == 0 {
			return nil, errors.New("no choices returned on follow-up")
		}

		finalMsg := finalResp.Choices[0].Message.Content

		// Save the interaction with final formatted response
		messageModel := &models.Message{
			UserID:      userID,
			UserContent: message,
			AssistantContent: finalMsg,
		}
		if _, err := s.iaSvc.StartOrAppendInteraction(ctx, 
			userID, 
			interactionID,
			messageModel,
			model,
			&prompt,
			); err != nil {
			return nil, err
		}

		return &models.BotResponse{
			Type:     "related_posts",
			Response: finalMsg,
		}, nil

	default:
		switch v := out.(type) {
		case string:
			return &models.BotResponse{Type: "string", Response: v}, nil
		case *models.Message:
			if _, err := s.iaSvc.StartOrAppendInteraction(ctx, 
				userID, 
				interactionID,
				v,
				model,
				&prompt,
				); err != nil {
				return nil, err
			}
			return &models.BotResponse{Type: "interaction", Response: v.AssistantContent}, nil
		case *models.Post:
			return &models.BotResponse{Type: "post", Response: v}, nil
		default:
			return nil, errors.New("unsupported function return type")
		}
	}
}

func (s *BotService) GenerateRequestStream(
	ctx context.Context,
	userID primitive.ObjectID,
	interactionID *primitive.ObjectID,
	// messages []openai.ChatCompletionMessageParamUnion,
	message string,
	systemPrompt *string,
	plugins []string,
	w http.ResponseWriter,
) {

	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = openai.ChatModelGPT4o
	}

	ctx = context.WithValue(ctx, bot.CtxUserID, userID)
	ctx = context.WithValue(ctx, bot.CtxInput, "")
	ctx = context.WithValue(ctx, bot.CtxRepo, s.pRepo)
	ctx = context.WithValue(ctx, bot.CtxClient, s.oaClient)

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	var prompt string
	var payloadStr string

	if  systemPrompt != nil {
		prompt = *systemPrompt
	} else {
		prompt = `
		Share a brief description of what you plan to do in 40 words.
		You are pretty good at whatever you are requested to do. 
		`
	}

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(prompt),
		openai.UserMessage(message),
	}


	var tools []openai.ChatCompletionToolParam
	if len(plugins) > 0 {
		tools = make([]openai.ChatCompletionToolParam, 0, len(bot.Registry))
		for _, spec := range bot.Registry {
			for _, p := range plugins {
				if p == spec.Definition.Name {
					tools = append(tools, openai.ChatCompletionToolParam{
						Function: openai.FunctionDefinitionParam{
							Name:        spec.Definition.Name,
							Description: openai.String(spec.Definition.Description),
							Parameters:  spec.Definition.Parameters,
						},
					})
				}
			}
		}
	}

	stream := s.oaClient.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Model:    model,
		Messages: messages,
		Tools:    tools,
		Seed:     openai.Int(0),
	})
	
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

		// will refactor to just reuse GenerateRequestStream

		if tool, ok := acc.JustFinishedToolCall(); ok {
			var toolStrOutput string

			// Check if tool needs a follow up
			if tool.Name == "get_related_posts" {
				spec, found := bot.Registry[tool.Name]
				if !found {
					http.Error(w, "unknown tool", http.StatusInternalServerError)
					return
				}
				
				// Parse arguments JSON into raw message
				args := json.RawMessage(tool.Arguments)
				
				out, err := spec.Handle(ctx, args)

				if err != nil {
					http.Error(w, fmt.Sprintf("tool handler error: %v", err), http.StatusInternalServerError)
					return
				}

				payloadBytes, ok := out.([]byte)
				if !ok {
					http.Error(w, "invalid tool output type", http.StatusInternalServerError)
					return
				}

				payloadStr = string(payloadBytes)
				prompt = `You just fetched related blog post data.

				Please format it as follows:

				1. Keep only these fields: slug, title, views, created_by.
				2. Start with a short, friendly intro (e.g. “Hey there! I found some posts you might like…”).
				3. Use **clean, well-formatted Markdown tables**. One table for relevant posts, another for not relevant.
				4. Use meaningful section headings like "### 🚀 Programming Posts" and "### 🌴 Other Interesting Reads".
				5. Keep the tone conversational but concise. Add a few light emojis for charm, but don’t overdo it.
				6. Keep spacing and formatting neat for maximum clarity.
				7. No explanations or bullet lists — only the intro and two tables.

				Thanks!`

			} else {
				spec, found := bot.Registry[tool.Name]
				if !found {
					http.Error(w, "unknown tool", http.StatusInternalServerError)
					return
				}
				
				// Parse arguments JSON into raw message
				args := json.RawMessage(tool.Arguments)
				out, err := spec.Handle(ctx, args)
				if err != nil {
					http.Error(w, fmt.Sprintf("tool handler error: %v", err), http.StatusInternalServerError)
					return
				}

				switch v := out.(type) {
				case []byte:
					payloadStr = string(v)
				case string:
					payloadStr = v
				default:
					// Try to marshal to JSON if it's a struct or something else
					b, err := json.Marshal(v)
					if err != nil {
						http.Error(w, "failed to marshal tool output", http.StatusInternalServerError)
						return
					}
					payloadStr = string(b)
				}

			}

			followupMessages := []openai.ChatCompletionMessageParamUnion{
				acc.Choices[0].Message.ToParam(),
				openai.ToolMessage(payloadStr, acc.Choices[0].Message.ToolCalls[0].ID),
				openai.SystemMessage(prompt),
				openai.UserMessage(message),
			}

			followupStream := s.oaClient.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
				Model:    model,
				Messages: followupMessages,
				Tools:    tools,
				Seed:     openai.Int(0),
			})

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

func escapeSSE(data string) string {
	data = strings.ReplaceAll(data, "\\", "\\\\")
	data = strings.ReplaceAll(data, "\n", "\\n")
	data = strings.ReplaceAll(data, "\r", "\\r")
	return data
}

