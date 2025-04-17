package services

import (
	"context"

	"github.com/AlexBetita/go_prac/internal/models"
	"github.com/AlexBetita/go_prac/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InteractionService struct {
	repo repositories.InteractionRepository
}

func NewInteractionService(repo repositories.InteractionRepository) *InteractionService {
	return &InteractionService{repo: repo}
}

func (s *InteractionService) SaveInteraction(
	ctx context.Context,
	userID primitive.ObjectID,
	userMsg string,
	botResp string,
) (*models.Interaction, error) {
	interaction := &models.Interaction{
		UserID:      userID,
		UserMessage: userMsg,
		BotResponse: botResp,
	}
	err := s.repo.Create(ctx, interaction)
	if err != nil {
		return nil, err
	}
	return interaction, nil
}
