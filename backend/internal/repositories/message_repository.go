package repositories

import (
	"context"
	"time"

	"github.com/AlexBetita/go_prac/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageRepository interface {
	Create(ctx context.Context, message *models.Message) error
	FindByInteractionID(ctx context.Context, interactionID primitive.ObjectID) ([]*models.Message, error)
	DeleteByInteractionIDs(ctx context.Context, ids []primitive.ObjectID) error
}

type mongoMessageRepository struct {
	coll *mongo.Collection
}

func NewMessageRepository(db *mongo.Database) MessageRepository {
	return &mongoMessageRepository{coll: db.Collection("messages")}
}

func (r *mongoMessageRepository) Create(ctx context.Context, message *models.Message) error {
	message.ID = primitive.NewObjectID()
	message.CreatedAt = time.Now().Unix()
	message.UpdatedAt = message.CreatedAt
	_, err := r.coll.InsertOne(ctx, message)
	return err
}

func (r *mongoMessageRepository) FindByInteractionID(ctx context.Context, interactionID primitive.ObjectID) ([]*models.Message, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"interaction_id": interactionID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []*models.Message
	for cursor.Next(ctx) {
		var msg models.Message
		if err := cursor.Decode(&msg); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}
	return messages, nil
}

func (r *mongoMessageRepository) DeleteByInteractionIDs(ctx context.Context, ids []primitive.ObjectID) error {
	if len(ids) == 0 {
		return nil
	}
	_, err := r.coll.DeleteMany(ctx, bson.M{"interaction_id": bson.M{"$in": ids}})
	return err
}
