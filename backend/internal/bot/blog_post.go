package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/AlexBetita/go_prac/internal/models"
	openai "github.com/openai/openai-go"
)

type BlogPostPayload struct {
	Content string   `json:"content"`
	Title   string   `json:"title"`
	Tags    []string `json:"tags"`
	Keywords    []string `json:"keywords"`
	Slug    string   `json:"slug"`
	Summary string   `json:"summary"`
	Will_Embed bool  `json:"will_embed"`
}

var blogPostDef = openai.FunctionDefinition{
	Name:        "create_blog_post",
	Description: "Return a markdown blog post about a given title",
	Parameters: map[string]any{
		"type": "object",
		"properties": map[string]any{
			"content": map[string]string{
				"type":        "string",
				"description": "Full blog post in markdown",
			},
			"title": map[string]string{
				"type": 	"string",
				"description": "Title of the blog post",
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
			"will_embed": map[string]any{
				"type": "boolean",
				"description": "True if user request to embed otherwise false",
				"default": false,
			},
		},
		"required": 
			[]string{"content", "title", 
			"tags", "summary", "keywords",
		"slug", "will_embed"},
	},
}

func blogPostHandler(ctx context.Context, raw json.RawMessage) (any, error) {

	var p BlogPostPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, err
	}

	post := &models.Post{
		UserID:    UserID(ctx),
		Message:   Input(ctx),
		Title:     p.Title,
		Content:   p.Content,
		Summary: p.Summary,
		Keywords: p.Keywords,
		Tags: p.Tags,
		Slug: generateSlug(p.Title),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	if p.Will_Embed {
		client := Client(ctx)
		embedding, err := EmbedText(ctx, client, p.Content)
		if err != nil {
			return nil, err
		}
		post.Embeddings = embedding.Data[0].Embedding
	}

	if err := Repo(ctx).Create(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func generateSlug(title string) string {
	base := slugify(title)
	suffix := rand.Intn(10000)
	return fmt.Sprintf("%s-%04d", base, suffix)
}

func slugify(s string) string {
	s = strings.ToLower(s)
	s = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

func init() {
	Register(Spec{
		Definition: blogPostDef,
		Handle:     blogPostHandler,
	})
}
