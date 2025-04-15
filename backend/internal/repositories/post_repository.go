package repositories

import (
	"context"
	"time"

	"github.com/AlexBetita/go_prac/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostRepository interface {
	Create(ctx context.Context, p *models.Post) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.Post, error)
}

type mongoPostRepository struct{ coll *mongo.Collection }

func NewPostRepository(db *mongo.Database) PostRepository {
	return &mongoPostRepository{coll: db.Collection("posts")}
}

func (r *mongoPostRepository) Create(ctx context.Context, p *models.Post) error {
	p.CreatedAt = time.Now().Unix()
	_, err := r.coll.InsertOne(ctx, p)
	return err
}

func (r *mongoPostRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Post, error) {
	var p models.Post
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&p)
	return &p, err
}