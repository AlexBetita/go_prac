package services

import (
	"context"
	"fmt"

	"github.com/AlexBetita/go_prac/internal/models"
	"github.com/AlexBetita/go_prac/internal/repositories"
	openai "github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/AlexBetita/go_prac/internal/bot"   
)

type PostService interface {
	GetPostsByID(ctx context.Context, id string) (*models.Post, error)
	SearchPosts(ctx context.Context, q string, limit int64) ([]*models.Post, error)
	SearchPostsByVector(ctx context.Context, q string, limit int64) ([]*models.Post, error)
}

type postService struct {
	repo      repositories.PostRepository
	oaClient *openai.Client
}

func NewPostService(repo repositories.PostRepository, oaClient *openai.Client) PostService {
	return &postService{repo: repo, oaClient: oaClient}
}

func (s *postService) GetPostsByID(ctx context.Context, id string) (*models.Post, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid post ID: %w", err)
	}

	post, err := s.repo.FindByID(ctx, objID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch post: %w", err)
	}

	return post, nil
}

func (s *postService) SearchPosts(ctx context.Context, q string, limit int64) ([]*models.Post, error) {
	return s.repo.Search(ctx, q, limit)
}

func (s *postService) SearchPostsByVector(ctx context.Context, q string, limit int64) ([]*models.Post, error) {
	vec, err := bot.EmbedText(ctx, s.oaClient, q)
	if err != nil {
		return nil, err
	}
	return s.repo.VectorSearch(ctx, vec, limit)
}
