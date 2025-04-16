package services

import (
	"context"
	"fmt"

	"github.com/AlexBetita/go_prac/internal/models"
	"github.com/AlexBetita/go_prac/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostService interface {
	GetBlogByID(ctx context.Context, id string) (*models.Post, error)
}

type postService struct {
	repo repositories.PostRepository
	jwtSecret string
}

func NewPostService(repo repositories.PostRepository, jwtSecret string) PostService {
	return &postService{repo: repo, jwtSecret: jwtSecret}
}

func (s *postService) GetBlogByID(ctx context.Context, id string) (*models.Post, error) {
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