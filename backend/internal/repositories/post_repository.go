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
    FindBySlug(ctx context.Context, slug string) (*models.Post, error)
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
        bson.D{{
            Key: "$search",
            Value: bson.D{
                {Key: "index", Value: "posts"},
                {Key: "compound", Value: bson.D{
                    {Key: "should", Value: bson.A{
                        bson.D{{
                            Key: "text",
                            Value: bson.D{
                                {Key: "query", Value: query},
                                {Key: "path",  Value: bson.A{"title", "content", "summary", "keywords", "tags"}},
                            },
                        }},
                    }},
                }},
            },
        }},
        bson.D{{Key: "$limit",   Value: limit}},
        bson.D{{Key: "$project", Value: bson.D{{Key: "embeddings", Value: 0}}}},
    }
    return r.aggregate(ctx, pipeline)
}

func (r *mongoPostRepository) VectorSearch(ctx context.Context, vector []float32, limit int64) ([]*models.Post, error) {
    pipeline := mongo.Pipeline{
        bson.D{{
            Key: "$vectorSearch",
            Value: bson.D{
                {Key: "index",         Value: "vector_index"},
                {Key: "path",          Value: "embeddings"},
                {Key: "queryVector",   Value: vector},
                {Key: "numCandidates", Value: limit * 5},
                {Key: "limit",         Value: limit},
            },
        }},
        bson.D{{
            Key: "$project",
            Value: bson.D{
                {Key: "embeddings", Value: 0},
            },
        }},
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

func (r *mongoPostRepository) FindBySlug(
    ctx context.Context,
    slug string,
) (*models.Post, error) {
    var p models.Post
    err := r.coll.FindOneAndUpdate(
        ctx,
        bson.M{"slug": slug},
        bson.M{"$inc": bson.M{"views": 1}},
        options.FindOneAndUpdate().SetReturnDocument(options.After),
    ).Decode(&p)
    return &p, err
}