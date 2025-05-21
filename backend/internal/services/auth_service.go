package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AlexBetita/go_prac/internal/models"
	"github.com/AlexBetita/go_prac/internal/repositories"
	"github.com/AlexBetita/go_prac/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthService interface {
    Register(ctx context.Context, email, password string) (*models.User, error)
    Login(ctx context.Context, email, password string) (string, *models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
}

type authService struct {
    repo repositories.UserRepository
    jwtSecret string
}

func NewAuthService(repo repositories.UserRepository, jwtSecret string) AuthService {
    return &authService{repo: repo, jwtSecret: jwtSecret}
}

func (s *authService) Register(ctx context.Context, email, password string) (*models.User, error) {
    user, err := s.repo.FindByEmail(ctx, email)
    if err == nil {
        return nil, errors.New("We couldn't create your account")
    }
    hashed, _ := utils.HashPassword(password)
    user = &models.User{
		Email:     email,
		Password:  hashed,
		Providers: []string{"local"},
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
    if err := s.repo.Create(ctx, user); err != nil {
        return nil, err
    }
    return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, *models.User, error) {
    user, err := s.repo.FindByEmail(ctx, email)
    if err != nil {
        return "", nil, err
    }
    if !utils.CheckPasswordHash(password, user.Password) {
        return "", nil, errors.New("invalid credentials")
    }
    token, err := utils.GenerateJWT(user.ID.Hex(), s.jwtSecret)
    return token, user, err
}

func (s *authService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, fmt.Errorf("invalid user ID: %w", err)
    }

    user, err := s.repo.FindByID(ctx, objID)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve user: %w", err)
    }

    return user, nil
}