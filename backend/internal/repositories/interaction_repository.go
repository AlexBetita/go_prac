package repositories

import (
	"context"
	"time"

	"github.com/AlexBetita/go_prac/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type InteractionRepository interface {
	Create(ctx context.Context, p *models.Interaction) error
	DeleteByFolderID(ctx context.Context, folderID primitive.ObjectID) error
	FindByFolderID(ctx context.Context, folderID primitive.ObjectID) ([]primitive.ObjectID, error)
}

type mongoInteractionRepository struct{ coll *mongo.Collection }

func NewInteractionRepository(db *mongo.Database) InteractionRepository {
	return &mongoInteractionRepository{coll: db.Collection("interactions")}
}

func (r *mongoInteractionRepository) Create(ctx context.Context, p *models.Interaction) error {
	p.CreatedAt = time.Now().Unix()
	res, err := r.coll.InsertOne(ctx, p)
	if err != nil {
		return err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		p.ID = oid
	}
	return nil
}

func (r *mongoInteractionRepository) DeleteByFolderID(ctx context.Context, folderID primitive.ObjectID) error {
	_, err := r.coll.DeleteMany(ctx, bson.M{"folder_id": folderID})
	return err
}

func (r *mongoInteractionRepository) FindByFolderID(ctx context.Context, folderID primitive.ObjectID) ([]primitive.ObjectID, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"folder_id": folderID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var ids []primitive.ObjectID
	for cursor.Next(ctx) {
		var interaction models.Interaction
		if err := cursor.Decode(&interaction); err != nil {
			return nil, err
		}
		ids = append(ids, interaction.ID)
	}
	return ids, nil
}