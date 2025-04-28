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

	var finalOut interface{}

	if len(msg.ToolCalls) == 0 {
		finalOut = &models.Interaction{
			UserID:      userID,
			UserMessage: message,
			BotResponse: msg.Content,
		}
	} else {
		call := msg.ToolCalls[0]
		spec, ok := bot.Registry[call.Function.Name]
		if !ok {
			return nil, errors.New("unknown function: " + call.Function.Name)
		}

		rawIn := json.RawMessage(call.Function.Arguments)
		out, err := spec.Handle(ctx, rawIn)
		if err != nil {
			return nil, err
		}

		finalOut = out

		if call.Function.Name == "get_related_posts" {
			payloadStr := string(out.([]byte))

			relatedPostsSystemMsg := openai.ChatCompletionMessage{
				Role: openai.ChatMessageRoleSystem,
				Content: `You just fetched related blog post data.

				Please format it as follows:

				1. Keep only these fields: slug, title, views, created_by.
				2. Start with a short, friendly intro (e.g. “Hey there! I found some posts you might like…”).
				3. Use **clean, well-formatted Markdown tables**. One table for relevant posts, another for not relevant.
				4. Use meaningful section headings like "### 🚀 Programming Posts" and "### 🌴 Other Interesting Reads".
				5. Keep the tone conversational but concise. Add a few light emojis for charm, but don’t overdo it.
				6. Keep spacing and formatting neat for maximum clarity.
				7. No explanations or bullet lists — only the intro and two tables.

				Thanks!`,
			}

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
				return nil, errors.New("no choices on follow-up call")
			}

			finalOut = finalResp.Choices[0].Message.Content
			finalStr, ok := finalOut.(string)
			if !ok {
				return nil, errors.New("finalOut is not a string")
			}
			saveRep := &models.Interaction{
				UserID:      userID,
				UserMessage: message,
				BotResponse: finalStr,
			}
			if err := s.iaRepo.Create(ctx, saveRep); err != nil {
				return nil, err
			}
		}
	}

	switch v := finalOut.(type) {
	case *models.Post:
		return &models.BotResponse{Type: "post", Response: v}, nil
	case *models.Interaction:
		if err := s.iaRepo.Create(ctx, v); err != nil {
			return nil, err
		}
		return &models.BotResponse{
			Type:     "interaction",
			Response: v.BotResponse,
		}, nil
	case string:
		return &models.BotResponse{
			Type:     "related_posts",
			Response: v,
		}, nil
	case *models.BotResponse:
		return v, nil
	default:
		return nil, errors.New("unsupported type returned")
	}
}
