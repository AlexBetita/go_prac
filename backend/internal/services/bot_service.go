package services

import (
	"context"
	"errors"
	"os"
	"encoding/json"

	"github.com/AlexBetita/go_prac/internal/bot"
	"github.com/AlexBetita/go_prac/internal/models"
	"github.com/AlexBetita/go_prac/internal/repositories"
	openai "github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BotService struct {
	repo   repositories.PostRepository
	client *openai.Client
}

func NewBotService(repo repositories.PostRepository, apiKey string) *BotService {
	return &BotService{repo: repo, client: openai.NewClient(apiKey)}
}

func (s *BotService) GenerateRequest(
	ctx context.Context,
	userID primitive.ObjectID,
	message string,
) (*models.Post, error) {

	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = openai.GPT4o20241120
	}

	ctx = context.WithValue(ctx, bot.CtxUserID, userID)
	ctx = context.WithValue(ctx, bot.CtxInput, message)
	ctx = context.WithValue(ctx, bot.CtxRepo, s.repo)

	tools := make([]openai.Tool, 0, len(bot.Registry))
	for _, spec := range bot.Registry {
		tools = append(tools, openai.Tool{
			Type:     openai.ToolTypeFunction,
			Function: &spec.Definition,
		})
	}

	req := openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: "You are pretty good at whatever you are requested to do."},
			{Role: openai.ChatMessageRoleUser,   Content: message},
		},
		Tools:      tools,
		ToolChoice: "auto",
	}

	resp, err := s.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, err
	}
	if len(resp.Choices) == 0 {
		return nil, errors.New("no choices returned by OpenAI")
	}

	msg := resp.Choices[0].Message
	if len(msg.ToolCalls) == 0 {
		return nil, errors.New("no tool call returned")
	}
	call := msg.ToolCalls[0]

	spec, ok := bot.Registry[call.Function.Name]
	if !ok {
		return nil, errors.New("unknown function: " + call.Function.Name)
	}

	raw := json.RawMessage(call.Function.Arguments)
	out, err := spec.Handle(ctx, raw)
	if err != nil {
		return nil, err
	}

	if post, ok := out.(*models.Post); ok {
		return post, nil
	}
	return nil, errors.New("handler did not return a Post")
}
