package config

import (
	"context"
	lc "lcgodemo/internal/langchain"
	"os"

	"github.com/tmc/langchaingo/llms"
)

type Provider int

const (
	providerUndefined Provider = iota
	google
	openai
	ollama
)

type errNoProviderError struct{}

func (e errNoProviderError) Error() string {
	return "no provider given on init"
}

func determineProvider() Provider {
	provider := os.Getenv("PROVIDER")

	if provider == "google" {
		return google
	}

	if provider == "openai" {
		return openai
	}

	if provider == "ollama" {
		return ollama
	}

	return providerUndefined
}

func parseApiKey() string {
	return os.Getenv("API_KEY")
}

func InitLLM(p Provider, ctx context.Context, keyOrModel string) (llms.Model, error) {
	switch p {
	case google:
		return lc.InitGoogleAI(ctx, keyOrModel)
	case openai:
		return lc.InitOpenAI(keyOrModel)
	case ollama:
		return lc.InitOllama(ctx, keyOrModel)
	default:
		return nil, new(errNoProviderError)
	}
}

func ParseEnv() (Provider, string) {
	provider := determineProvider()
	apiKey := parseApiKey()
	return provider, apiKey
}
