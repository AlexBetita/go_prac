package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Email     string             `bson:"email" json:"email"`
    Password  string             `bson:"password,omitempty" json:"-"`
    Provider  string             `bson:"provider" json:"provider"`
    ProviderID string            `bson:"provider_id,omitempty" json:"provider_id,omitempty"`
    CreatedAt int64              `bson:"created_at" json:"created_at"`
    UpdatedAt int64              `bson:"updated_at" json:"updated_at"`
}