package lc

import (
	"context"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

// InitOllama is the entrypoint for interacting with Ollama provided
// LLMs
func InitOllama(ctx context.Context, model string) (llms.Model, error) {
	const serverurl = "http://localhost:11434/api/generate"
	llm, err := ollama.New(ollama.WithModel(model), ollama.WithServerURL(serverurl))
	if err != nil {
		return nil, err
	}

	return llm, nil
}
