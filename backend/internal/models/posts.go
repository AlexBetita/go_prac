package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Post struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty"    json:"-"`
	UserID      primitive.ObjectID   `bson:"user_id"          json:"-"`
	Title       string               `bson:"title"            json:"title"`
	Content     string               `bson:"content"          json:"content"`
	Summary     string               `bson:"summary"          json:"summary"`
	Message     string               `bson:"message"          json:"-"`
	Keywords    []string             `bson:"keywords"         json:"keywords"`
	Tags        []string             `bson:"tags"             json:"tags"`
	Slug        string               `bson:"slug"             json:"slug"`
	Views       int64                `bson:"views"            json:"views"`
	Embeddings  []float32            `bson:"embeddings"       json:"-"`
	CreatedBy   string				 `bson:"created_by"       json:"created_by"`
	CreatedAt   int64                `bson:"created_at"       json:"created_at"`
	UpdatedAt   int64                `bson:"updated_at"       json:"-"`
}