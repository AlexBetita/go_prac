package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
    ID        primitive.ObjectID    `bson:"_id,omitempty" json:"-"`
    Email     string                `bson:"email" json:"email"`
    Password  string                `bson:"password,omitempty" json:"-"`
    Providers  []string             `bson:"providers" json:"providers"`
    ProviderID string               `bson:"provider_id,omitempty" json:"-"`
    Role       string               `bson:"role,omitempty" json:"role,omitempty"`
    APIKeys    map[string]string    `bson:"api_keys" json:"api_keys"` // { "OpenAI": "KEY", "Anthropic": "KEY" }
    LastActive int64                `bson:"last_active" json:"last_active"`
    CreatedAt int64                 `bson:"created_at" json:"created_at"`
    UpdatedAt int64                 `bson:"updated_at" json:"updated_at"`
}