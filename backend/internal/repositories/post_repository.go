package repositories

import (
	"context"
	"time"

	"github.com/AlexBetita/go_prac/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostRepository interface {
	Create(ctx context.Context, p *models.Post) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.Post, error)
	Search(ctx context.Context, query string, limit int64) ([]*models.Post, error)
	VectorSearch(ctx context.Context, vector []float32, limit int64) ([]*models.Post, error)
}

type mongoPostRepository struct{ coll *mongo.Collection }

func NewPostRepository(db *mongo.Database) PostRepository {
	return &mongoPostRepository{coll: db.Collection("posts")}
}

func (r *mongoPostRepository) Create(ctx context.Context, p *models.Post) error {
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

func (r *mongoPostRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Post, error) {
	var p models.Post
	err := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"_id": id},
		bson.M{"$inc": bson.M{"views": 1}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&p)
	return &p, err
}

func (r *mongoPostRepository) Search(ctx context.Context, query string, limit int64) ([]*models.Post, error) {
    pipeline := mongo.Pipeline{

        bson.D{{"$search", bson.D{
            {"index", "posts"},
            {"compound", bson.D{
                {"should", bson.A{
                    bson.D{{"text", bson.D{
                        {"query", query},
                        {"path", bson.A{"topic", "content", "summary", "keywords", "tags"}},
                    }}},
                }},
            }},
        }}},

        bson.D{{"$limit", limit}},

        bson.D{{"$project", bson.D{{"embeddings", 0}}}},
    }
    return r.aggregate(ctx, pipeline)
}

func (r *mongoPostRepository) VectorSearch(ctx context.Context, vector []float32, limit int64) ([]*models.Post, error) {
    pipeline := mongo.Pipeline{
        bson.D{{"$search", bson.D{
            {"index", "vector_index"},
            {"knn", bson.D{
                {"vector", vector},
                {"path",   "embeddings"},
                {"k",      limit},
            }},
        }}},
        bson.D{{"$project", bson.D{{"embeddings", 0}}}},
    }
    return r.aggregate(ctx, pipeline)
}

func (r *mongoPostRepository) aggregate(ctx context.Context, pipeline mongo.Pipeline) ([]*models.Post, error) {
    cur, err := r.coll.Aggregate(ctx, pipeline)
    if err != nil {
        return nil, err
    }
    defer cur.Close(ctx)

    var out []*models.Post
    if err := cur.All(ctx, &out); err != nil {
        return nil, err
    }
    return out, nil
}