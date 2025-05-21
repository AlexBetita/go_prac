package repositories

import (
	"context"
	"time"

	"github.com/AlexBetita/go_prac/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FolderRepository interface {
	Create(ctx context.Context, folder *models.Folder) error
	Update(ctx context.Context, id primitive.ObjectID, update bson.M) error
	FindMany(ctx context.Context, filter bson.M) ([]*models.Folder, error)
	FindOne(ctx context.Context, id primitive.ObjectID) (*models.Folder, error)
}

type mongoFolderRepo struct {
	coll *mongo.Collection
}

func NewFolderRepository(db *mongo.Database) FolderRepository {
	return &mongoFolderRepo{coll: db.Collection("folders")}
}

func (r *mongoFolderRepo) Create(ctx context.Context, folder *models.Folder) error {
	folder.CreatedAt = time.Now().Unix()
	folder.UpdatedAt = folder.CreatedAt
	res, err := r.coll.InsertOne(ctx, folder)
	if err == nil {
		folder.ID = res.InsertedID.(primitive.ObjectID)
	}
	return err
}

func (r *mongoFolderRepo) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	update["updated_at"] = time.Now().Unix()
	_, err := r.coll.UpdateByID(ctx, id, bson.M{"$set": update})
	return err
}

func (r *mongoFolderRepo) FindMany(ctx context.Context, filter bson.M) ([]*models.Folder, error) {
	cursor, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var folders []*models.Folder
	for cursor.Next(ctx) {
		var folder models.Folder
		if err := cursor.Decode(&folder); err != nil {
			return nil, err
		}
		folders = append(folders, &folder)
	}
	return folders, nil
}

func (r *mongoFolderRepo) FindOne(ctx context.Context, id primitive.ObjectID) (*models.Folder, error) {
	var folder models.Folder
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&folder)
	return &folder, err
}
