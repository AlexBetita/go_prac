package bot

import (
	"context"
	"fmt"
	
	openai "github.com/openai/openai-go"
)

func EmbedText(ctx context.Context, oaClient *openai.Client, input string) (*openai.CreateEmbeddingResponse, error) {
	resp, err := oaClient.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String(input),
		},
		Model: openai.EmbeddingModelTextEmbeddingAda002,
	})
	if err != nil {
        return nil, fmt.Errorf("failed to embed input: %w", err)
    }
    if len(resp.Data) == 0 {
        return nil, fmt.Errorf("no embedding returned")
    }
    return resp, nil
}