package lc

import (
	"context"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

func getOllamaServerAddres() string {
	if os.Getenv("OLLAMA_ADDRESS") != "" {
		return os.Getenv("OLLAMA_ADDRESS")
	}

	return "localhost"
}

// InitOllama is the entrypoint for interacting with Ollama provided
// LLMs
func InitOllama(ctx context.Context, model string) (llms.Model, error) {
	var serverurl = "http://" + getOllamaServerAddres() + ":11434"
	llm, err := ollama.New(ollama.WithModel(model), ollama.WithServerURL(serverurl))
	if err != nil {
		return nil, err
	}

	return llm, nil
}
