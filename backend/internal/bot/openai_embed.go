package bot

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

func EmbedText(ctx context.Context, oaClient *openai.Client, input string) ([]float32, error) {
	resp, err := oaClient.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: []string{input},
		Model: openai.AdaEmbeddingV2,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to embed input: %w", err)
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}
	return resp.Data[0].Embedding, nil
}
