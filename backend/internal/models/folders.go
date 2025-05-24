package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Folder struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FolderName      string             `bson:"folder_name" json:"folder_name"`
	FolderContext   string             `bson:"folder_context" json:"folder_context"`
	FolderDocuments []string           `bson:"folder_documents" json:"folder_documents"`
	Metadata     map[string]interface{} `bson:"metadata" json:"metadata"`
	CreatedBy       primitive.ObjectID `bson:"created_by" json:"created_by"`
	Favorite        bool               `bson:"favorite" json:"favorite"`
	CreatedAt       int64              `bson:"created_at" json:"created_at"`
	UpdatedAt       int64              `bson:"updated_at" json:"updated_at"`
}
