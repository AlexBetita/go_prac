package services

import (
	"context"
	"fmt"
	"time"

	"github.com/AlexBetita/go_prac/internal/models"
	"github.com/AlexBetita/go_prac/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InteractionService struct {
	repo repositories.InteractionRepository
	messageRepo  repositories.MessageRepository
}

func NewInteractionService(repo repositories.InteractionRepository, msgRepo repositories.MessageRepository) *InteractionService {
	return &InteractionService{repo: repo, messageRepo: msgRepo}
}

func (s *InteractionService) StartOrAppendInteraction(
	ctx context.Context,
	userID primitive.ObjectID,
	interactionID *primitive.ObjectID,
	message *models.Message,
	model string,
	systemPrompt *string,
) (*models.Interaction, error) {
	var interaction *models.Interaction
	var err error

	if interactionID == nil {
		// New chat
		newInteraction := &models.Interaction{
			UserID:      userID,
			Title:       "New Chat",
			Tags:        []string{},
			Favorite:    false,
			Metadata:    map[string]interface{}{},
			DefaultModel: model,
			CurrentModel: model,
			SystemPrompt: func() string {
			if systemPrompt != nil { return *systemPrompt }
				return "You are pretty good at whatever you are requested to do."
			}(),
			CreatedAt:   time.Now().Unix(),
			UpdatedAt:   time.Now().Unix(),
		}
		if err = s.repo.Create(ctx, newInteraction); err != nil {
			return nil, err
		}
		interaction = newInteraction
	} else {
		interaction, err = s.repo.GetByID(ctx, *interactionID)
		if err != nil {
			return nil, err
		}
	}

	// Add message
	message.InteractionID = interaction.ID
	message.UserID = userID
	message.CreatedAt = time.Now().Unix()
	message.UpdatedAt = message.CreatedAt
	if err := s.messageRepo.Create(ctx, message); err != nil {
		return nil, err
	}
	return interaction, nil
}

func (s *InteractionService) UpdateInteraction(
	ctx context.Context,
	id primitive.ObjectID,
	updates map[string]interface{},
) error {
	if val, ok := updates["folder_id"]; ok {
		strID, ok := val.(string)
		if !ok {
			return fmt.Errorf("folder_id must be a string")
		}

		objID, err := primitive.ObjectIDFromHex(strID)
		if err != nil {
			return fmt.Errorf("invalid folder_id format: %w", err)
		}

		updates["folder_id"] = objID
	}
	return s.repo.Update(ctx, id, updates)
}

