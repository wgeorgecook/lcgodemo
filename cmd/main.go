package main

import (
	"context"
	"log"

	"lcgodemo/internal/config"
	"lcgodemo/internal/rag"

	"github.com/tmc/langchaingo/llms"
)

func main() {
	log.Println("Hello!")
	defer log.Println("goodbye!	")
	provider, apiKey := config.ParseEnv()
	llm, err := config.InitLLM(provider, context.TODO(), apiKey)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	prompt := "What would be a good company name for a company that makes colorful socks? " + rag.LoadGroundingContext()
	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(completion)
}
