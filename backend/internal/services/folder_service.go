package services

import (
	"context"
	"errors"
	"time"

	"github.com/AlexBetita/go_prac/internal/models"
	"github.com/AlexBetita/go_prac/internal/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FolderService struct {
	repo repositories.FolderRepository
	interactionRepo repositories.InteractionRepository
	messageRepo repositories.MessageRepository
}

func NewFolderService(folderRepo repositories.FolderRepository, interactionRepo repositories.InteractionRepository,
	messageRepo repositories.MessageRepository) *FolderService {
	return &FolderService{
		repo:            folderRepo,
		interactionRepo: interactionRepo,
		messageRepo: messageRepo,
	}
}

func (s *FolderService) CreateFolder(ctx context.Context, folder *models.Folder) error {
	return s.repo.Create(ctx, folder)
}

func (s *FolderService) UpdateFolder(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	update["updated_at"] = time.Now().Unix()
	return s.repo.Update(ctx, id, update)
}

func (s *FolderService) GetFavorites(ctx context.Context, userID primitive.ObjectID) ([]*models.Folder, error) {
	return s.repo.FindMany(ctx, bson.M{"created_by": userID, "favorite": true})
}

func (s *FolderService) GetAllFolders(ctx context.Context, userID primitive.ObjectID) ([]*models.Folder, error) {
	return s.repo.FindMany(ctx, bson.M{"created_by": userID})
}

func (s *FolderService) GetFolder(ctx context.Context, id primitive.ObjectID) (*models.Folder, error) {
	return s.repo.FindOne(ctx, id)
}

func (s *FolderService) GetFoldersPaginated(ctx context.Context, userID primitive.ObjectID, limit, skip int) ([]*models.Folder, int64, error) {
	return s.repo.FindManyPaginated(ctx, bson.M{"created_by": userID}, limit, skip)
}

func (s *FolderService) DeleteFolder(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.Delete(ctx, id)
}

func (s *FolderService) DeleteInteractionsByFolder(ctx context.Context, folderID primitive.ObjectID) error {
	if s.interactionRepo == nil || s.messageRepo == nil {
		return errors.New("repositories not initialized")
	}

	// Get interaction IDs in this folder
	interactionIDs, err := s.interactionRepo.FindByFolderID(ctx, folderID)
	if err != nil {
		return err
	}

	// Delete messages
	if err := s.messageRepo.DeleteByInteractionIDs(ctx, interactionIDs); err != nil {
		return err
	}

	// Delete interactions
	return s.interactionRepo.DeleteByFolderID(ctx, folderID)
}

func (s *FolderService) ToggleFavorite(ctx context.Context, id primitive.ObjectID, value bool) error {
	return s.repo.Update(ctx, id, bson.M{"favorite": value})
}

func (s *FolderService) GetFavoriteFolders(ctx context.Context, userID primitive.ObjectID) ([]*models.Folder, error) {
	return s.repo.FindFavoritesByUser(ctx, userID)
}

func (s *FolderService) GetFavoriteFoldersPaginated(ctx context.Context, userID primitive.ObjectID, limit, skip int) ([]*models.Folder, int64, error) {
	return s.repo.FindFavoritePaginated(ctx, userID, limit, skip)
}
