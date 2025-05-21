package services

import (
	"context"
	"time"

	"github.com/AlexBetita/go_prac/internal/models"
	"github.com/AlexBetita/go_prac/internal/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FolderService struct {
	repo repositories.FolderRepository
}

func NewFolderService(repo repositories.FolderRepository) *FolderService {
	return &FolderService{repo: repo}
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
