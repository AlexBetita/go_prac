package bot

import (
	"context"

	"github.com/AlexBetita/go_prac/internal/repositories"
	openai "github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ctxKey string

const (
	CtxUserID ctxKey = "userID"
	CtxInput  ctxKey = "userInput"
	CtxRepo   ctxKey = "postRepo"
	CtxClient ctxKey = "openaiClient"
)

func UserID(ctx context.Context) primitive.ObjectID {
	return ctx.Value(CtxUserID).(primitive.ObjectID)
}
func Input(ctx context.Context) string { return ctx.Value(CtxInput).(string) }
func Repo(ctx context.Context) repositories.PostRepository {
	return ctx.Value(CtxRepo).(repositories.PostRepository)
}
func Client(ctx context.Context) *openai.Client {
	return ctx.Value(CtxClient).(*openai.Client)
}