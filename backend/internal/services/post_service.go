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
	GetPost(ctx context.Context, identifier string) (*models.Post, error)
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

func (s *postService) GetPost(ctx context.Context, identifier string) (*models.Post, error) {
	if oid, err := primitive.ObjectIDFromHex(identifier); err == nil {
		if post, err := s.repo.FindByID(ctx, oid); err == nil {
			return post, nil
		} else {
			return nil, fmt.Errorf("failed to find post by ID %s: %w", oid.Hex(), err)
		}
	}
	
	post, err := s.repo.FindBySlug(ctx, identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to find post by slug '%s': %w", identifier, err)
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
