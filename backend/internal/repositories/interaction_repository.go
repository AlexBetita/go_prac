package repositories

import (
	"context"
	"time"

	"github.com/AlexBetita/go_prac/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type InteractionRepository interface {
	Create(ctx context.Context, p *models.Interaction) error
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