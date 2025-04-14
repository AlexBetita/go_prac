package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/AlexBetita/go_prac/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    FindByEmail(ctx context.Context, email string) (*models.User, error)
    FindByID(ctx context.Context, id primitive.ObjectID) (*models.User, error)
}

type mongoUserRepository struct {
    collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
    return &mongoUserRepository{collection: db.Collection("users")}
}

func (r *mongoUserRepository) Create(ctx context.Context, user *models.User) error {
    user.CreatedAt = time.Now().Unix()
    user.UpdatedAt = user.CreatedAt
    _, err := r.collection.InsertOne(ctx, user)
    return err
}

func (r *mongoUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
    var user models.User
    err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
    if err == mongo.ErrNoDocuments {
        return nil, errors.New("user not found")
    }
    return &user, err
}

func (r *mongoUserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
    var user models.User
    err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
    if err == mongo.ErrNoDocuments {
        return nil, errors.New("user not found")
    }
    return &user, err
}
