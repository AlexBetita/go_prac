package bot

import (
	"context"
	"encoding/json"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

type RelatedPostsPayload struct {
	PostContent string `json:"post_content"`
	Count       int    `json:"count"`
}

type SlimPost struct {
    Slug      string `json:"slug"`
    Topic     string `json:"topic"`
    Views     int64    `json:"views"`
    CreatedBy string `json:"created_by"`
}

var relatedPostsDef = openai.FunctionDefinition{
	Name:        "get_related_posts",
	Description: `Get the top N blog posts, posts or blogs related to the provided content.`,
	Parameters: map[string]any{
		"type": "object",
		"properties": map[string]any{
			"post_content": map[string]string{
				"type":        "string",
				"description": "The content to find related blog posts, posts, or blogs for",
			},
			"count": map[string]any{
                "type":        "integer",
                "description": "How many related posts to return",
                "default":     5,
            },
		},
		"required": []string{"post_content"},
	},
}

func relatedPostsHandler(ctx context.Context, raw json.RawMessage) (any, error) {
	var p RelatedPostsPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, fmt.Errorf("invalid payload: %w", err)
	}
	if p.Count <= 0 {
        p.Count = 5
    }
	client := Client(ctx)
	repo := Repo(ctx)

	embedding, err := EmbedText(ctx, client, p.PostContent)
	if err != nil {
		return nil, fmt.Errorf("embedding failed: %w", err)
	}

	posts, err := repo.VectorSearch(ctx, embedding, int64(p.Count))
	if err != nil {
		return nil, fmt.Errorf("vector search failed: %w", err)
	}

	slims := make([]SlimPost, len(posts))
	for i, p := range posts {
		slims[i] = SlimPost{
			Slug:      p.Slug,
			Topic:     p.Topic,
			Views:     p.Views,
			CreatedBy: p.CreatedBy,
		}
	}
	payload, err := json.Marshal(slims)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func init() {
	Register(Spec{
		Definition: relatedPostsDef,
		Handle:     relatedPostsHandler,
	})
}
