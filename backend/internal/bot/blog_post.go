package bot

import (
	"context"
	"encoding/json"
	"time"

	openai "github.com/sashabaranov/go-openai"
	"github.com/AlexBetita/go_prac/internal/models"
)

var blogPostDef = openai.FunctionDefinition{
	Name:        "create_blog_post",
	Description: "Return a markdown blog post about a given topic",
	Parameters: map[string]any{
		"type": "object",
		"properties": map[string]any{
			"content": map[string]string{
				"type":        "string",
				"description": "Full blog post in markdown",
			},
		},
		"required": []string{"content"},
	},
}

func blogPostHandler(ctx context.Context, raw json.RawMessage) (any, error) {

	var p struct{ Content string `json:"content"` }
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, err
	}

	post := &models.Post{
		UserID:    UserID(ctx),
		Topic:     Input(ctx),
		Content:   p.Content,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	if err := Repo(ctx).Create(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func init() {
	Register(Spec{
		Definition: blogPostDef,
		Handle:     blogPostHandler,
	})
}
