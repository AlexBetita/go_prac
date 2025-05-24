package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/AlexBetita/go_prac/internal/bot"
	"github.com/AlexBetita/go_prac/internal/models"
	"github.com/AlexBetita/go_prac/internal/repositories"
	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/ssestream"
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
	messages []openai.ChatCompletionMessageParamUnion,
	plugins []string,
) *ssestream.Stream[openai.ChatCompletionChunk] {

	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = openai.ChatModelGPT4o
	}

	ctx = context.WithValue(ctx, bot.CtxUserID, userID)
	ctx = context.WithValue(ctx, bot.CtxInput, "")
	ctx = context.WithValue(ctx, bot.CtxRepo, s.pRepo)
	ctx = context.WithValue(ctx, bot.CtxClient, s.oaClient)

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
	

	return stream
}

