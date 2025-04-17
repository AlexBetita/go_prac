package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Interaction struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty"    json:"id"`
	UserID      primitive.ObjectID   `bson:"user_id"          json:"user_id"`
	UserMessage string         	     `bson:"user_message"     json:"user_message"`
	BotResponse string               `bson:"bot_response"     json:"bot_response"`
	CreatedAt   int64                `bson:"created_at"       json:"created_at"`
}