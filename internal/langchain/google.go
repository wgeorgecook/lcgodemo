package lc

import (
	"context"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
)

// InitGoogleAI is the entrypoint for interacting with Google LLMs
func InitGoogleAI(ctx context.Context, apiKey string) (llms.Model, error) {
	client, err := googleai.New(ctx, googleai.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	return client, nil
}
