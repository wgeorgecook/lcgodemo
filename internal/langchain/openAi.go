package lc

import (
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

// Init is the entrypoint for interacting with OpenAI LLMs
func InitOpenAI(apiKey string) (llms.Model, error) {
	llm, err := openai.New(openai.WithToken(apiKey))
	if err != nil {
		return nil, err
	}

	return llm, nil
}
