package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Interaction struct {
	ID          primitive.ObjectID   	`bson:"_id,omitempty"    json:"id"`
	UserID      primitive.ObjectID   	`bson:"user_id"          json:"-"`
	FolderID    *primitive.ObjectID  	`bson:"folder_id,omitempty" json:"folder_id,omitempty"`
	Title        string              	`bson:"title" json:"title"`
	Tags         []string            	`bson:"tags" json:"tags"`
	SystemPrompt string              	`bson:"system_prompt, omitempty" json:"system_prompt,omitempty"`
	DefaultModel    string              `bson:"default_model" json:"default_model"`
	CurrentModel    string              `bson:"current_model" json:"current_model"`
	Metadata     map[string]interface{} `bson:"metadata" json:"metadata"`
	Favorite     bool                   `bson:"favorite" json:"favorite"`
	CreatedAt    int64                  `bson:"created_at" json:"created_at"`
	UpdatedAt    int64                  `bson:"updated_at" json:"updated_at"`
}