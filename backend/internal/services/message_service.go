package services

import (
	"context"

	"github.com/AlexBetita/go_prac/internal/models"
	"github.com/AlexBetita/go_prac/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageService struct {
	repo repositories.MessageRepository
}

func NewMessageService(r repositories.MessageRepository) *MessageService {
	return &MessageService{repo: r}
}

func (s *MessageService) Create(ctx context.Context, msg *models.Message) error {
	return s.repo.Create(ctx, msg)
}

func (s *MessageService) GetByInteraction(ctx context.Context, interactionID primitive.ObjectID) ([]*models.Message, error) {
	return s.repo.FindByInteractionID(ctx, interactionID)
}
