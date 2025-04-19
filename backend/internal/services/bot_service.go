package services

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/AlexBetita/go_prac/internal/bot"
	"github.com/AlexBetita/go_prac/internal/models"
	"github.com/AlexBetita/go_prac/internal/repositories"
	openai "github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BotService struct {
	pRepo   repositories.PostRepository
	iaRepo repositories.InteractionRepository
	oaClient *openai.Client
}

func NewBotService(pRepo repositories.PostRepository,
	iaRepo repositories.InteractionRepository,
	oaClient *openai.Client) *BotService {
	return &BotService{pRepo: pRepo, iaRepo: iaRepo, oaClient: oaClient}
}

func (s *BotService) GenerateRequest(
	ctx context.Context,
	userID primitive.ObjectID,
	message string,
) (*models.BotResponse, error) {

	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = openai.GPT4o20241120
	}

	systemMsg := openai.ChatCompletionMessage{
        Role:    openai.ChatMessageRoleSystem,
        Content: "You are pretty good at whatever you are requested to do.",
    }
    userMsg := openai.ChatCompletionMessage{
        Role:    openai.ChatMessageRoleUser,
        Content: message,
    }

	ctx = context.WithValue(ctx, bot.CtxUserID, userID)
	ctx = context.WithValue(ctx, bot.CtxInput, message)
	ctx = context.WithValue(ctx, bot.CtxRepo, s.pRepo)
	ctx = context.WithValue(ctx, bot.CtxClient, s.oaClient)

	tools := make([]openai.Tool, 0, len(bot.Registry))
	for _, spec := range bot.Registry {
		tools = append(tools, openai.Tool{
			Type:     openai.ToolTypeFunction,
			Function: &spec.Definition,
		})
	}

	firstReq := openai.ChatCompletionRequest{
        Model:      model,
        Messages:   []openai.ChatCompletionMessage{systemMsg, userMsg},
        Tools:      tools,
        ToolChoice: "auto",
    }

    firstResp, err := s.oaClient.CreateChatCompletion(ctx, firstReq)
    if err != nil {
        return nil, err
    }

    if len(firstResp.Choices) == 0 {
        return nil, errors.New("no choices returned by OpenAI")
    }

    msg := firstResp.Choices[0].Message

	if len(msg.ToolCalls) == 0 {
		interaction := &models.Interaction{
			UserID:      userID,
			UserMessage: message,
			BotResponse: msg.Content,
		}
		if err := s.iaRepo.Create(ctx, interaction); err != nil {
			return nil, err
		}
		return &models.BotResponse{
			Type:     "interaction",
			Response: interaction,
		}, nil
	}

	call := msg.ToolCalls[0]
	spec, ok := bot.Registry[call.Function.Name]
	if !ok {
		return nil, errors.New("unknown function: " + call.Function.Name)
	}

	if call.Function.Name == "get_related_posts" {
        rawIn := json.RawMessage(call.Function.Arguments)
		relatedPostsSystemMsg := openai.ChatCompletionMessage{
    		Role: openai.ChatMessageRoleSystem,
    		Content: `You’ve just fetched some related blog‑post data. 
			Please:
			1. Keep only these fields: slug, topic, views, created_by.  
			2. Write a brief, conversational intro (“Hey there! I found these related posts…”).  
			3. Render the posts in a Markdown list or table (your choice), showing only those four fields.  
			4. Give it your usual friendly flair—feel free to sprinkle in a couple of emojis or asides to keep it engaging.
			5. Seperate the lists by relevant to not relevant`,
		}
        resultJSON, err := spec.Handle(ctx, rawIn)
        if err != nil {
            return nil, err
        }
		payloadStr := string(resultJSON.([]byte))
        followupReq := openai.ChatCompletionRequest{
            Model: model,
            Messages: []openai.ChatCompletionMessage{
                relatedPostsSystemMsg,
                userMsg,
                {
                    Role:    openai.ChatMessageRoleFunction,
                    Name:    call.Function.Name,
                    Content: payloadStr,
                },
            },
        }
        finalResp, err := s.oaClient.CreateChatCompletion(ctx, followupReq)
        if err != nil {
            return nil, err
        }
        if len(finalResp.Choices) == 0 {
            return nil, errors.New("no choices on follow‑up call")
        }
        return &models.BotResponse{
            Type:     "related_posts",
            Response: finalResp.Choices[0].Message.Content,
        }, nil
    }

	raw := json.RawMessage(call.Function.Arguments)
	out, err := spec.Handle(ctx, raw)
	if err != nil {
		return nil, err
	}

	switch v := out.(type) {
		case *models.Post:
			return &models.BotResponse{Type: "post", Response: v}, nil
		case []*models.Post:
			return &models.BotResponse{Type: "related_posts", Response: v}, nil
		case *models.Interaction:
			return &models.BotResponse{Type: "interaction", Response: v}, nil
		case *models.BotResponse:
			return v, nil
		default:
			return nil, errors.New("tool `" + call.Function.Name + "` returned unsupported type")
    }
}
