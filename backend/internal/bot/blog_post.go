package bot

import (
	"context"
	"encoding/json"
	"time"

	openai "github.com/sashabaranov/go-openai"
	"github.com/AlexBetita/go_prac/internal/models"
)

type BlogPostPayload struct {
	Content string   `json:"content"`
	Topic   string   `json:"topic"`
	Tags    []string `json:"tags"`
	Keywords    []string `json:"keywords"`
	Slug    string   `json:"slug"`
	Summary string   `json:"summary"`
}

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
			"topic": map[string]string{
				"type": 	"string",
				"description": "Topic of the blog post",
			},
			"tags": map[string]any{
				"type": "array",
				"items": map[string]string{"type": "string"},
				"description": `Relevant tags for the blog post. Use this list [Tech, Life, News, 
				Education, Health, Entertainment, Business, Travel, Science, Art]`,
			},
			"summary": map[string]string{
				"type": "string",
				"description": "A short summary of the blog post. 1 or 2 sentences.",
			},
			"keywords": map[string]any{
				"type": "array",
				"items": map[string]string{"type": "string"},
				"description": "Curate relevant keywords.",
			},
			"slug": map[string]string{
				"type": "string",
				"description": "Create an SEO friendly slug",
			},
		},
		"required": 
			[]string{"content", "topic", 
			"tags", "summary", "keywords",
		"slug"},
	},
}

func blogPostHandler(ctx context.Context, raw json.RawMessage) (any, error) {

	var p BlogPostPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, err
	}

	post := &models.Post{
		UserID:    UserID(ctx),
		Topic:     p.Topic,
		Content:   p.Content,
		Summary: p.Summary,
		Keywords: p.Keywords,
		Tags: p.Tags,
		Slug: p.Slug,
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
