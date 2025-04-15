package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Post struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id"       json:"user_id"`
	Topic     string             `bson:"topic"         json:"topic"`
	Content   string             `bson:"content"       json:"content"`
	CreatedAt int64              `bson:"created_at"    json:"created_at"`
	UpdatedAt int64              `bson:"updated_at"    json:"updated_at"`
}