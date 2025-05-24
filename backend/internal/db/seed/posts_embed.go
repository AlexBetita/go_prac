package seed

import (
	"context"
	"log"

	"github.com/AlexBetita/go_prac/internal/bot"
	openai "github.com/openai/openai-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func EmbedSeededPosts(ctx context.Context, db *mongo.Database, client *openai.Client) error {
	collection := db.Collection("posts")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var post struct {
			ID      any    `bson:"_id"`
			Content string `bson:"content"`
		}
		if err := cursor.Decode(&post); err != nil {
			log.Println("⚠️ Failed to decode post:", err)
			continue
		}

		embedding, err := bot.EmbedText(ctx, client, post.Content)
		if err != nil {
			log.Println("⚠️ Failed to embed post:", post.ID, err)
			continue
		}

		_, err = collection.UpdateByID(ctx, post.ID, bson.M{"$set": bson.M{"embeddings": embedding}})
		if err != nil {
			log.Println("⚠️ Failed to update post:", post.ID, err)
		} else {
			log.Println("✅ Embedded:", post.ID)
		}
	}

	return cursor.Err()
}
