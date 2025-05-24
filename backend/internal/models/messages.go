package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Message struct {
	ID                        primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	UserID                    primitive.ObjectID     `bson:"user_id" json:"user_id"`
	InteractionID             primitive.ObjectID     `bson:"interaction_id" json:"interaction_id"`
	UserContent               string                 `bson:"user_content" json:"user_content"`
	AssistantContent          string                 `bson:"assistant_content" json:"assistant_content"`
	Files					  []string				 `bson:"files" json:"files"`
	Metadata                  map[string]interface{} `bson:"metadata" json:"metadata"` // { model: { context_size, context_size_with_sys_prompt } }
	CreatedAt                 int64                  `bson:"created_at" json:"created_at"`
	UpdatedAt                 int64                  `bson:"updated_at" json:"updated_at"`
}

type Attachment struct {
    Name    string `json:"name"`
    URL     string `json:"url"`
    Type    string `json:"type"`
    IsImage bool   `json:"is_image"`
}
