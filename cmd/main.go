package main

import (
	"context"
	"log"
	"os"

	"lcgodemo/internal/config"
	"lcgodemo/internal/rag"
	"lcgodemo/internal/vectors"

	"github.com/tmc/langchaingo/llms"
)

func main() {

	// set up the application
	log.Println("Hello!")
	defer log.Println("goodbye!	")
	provider, apiKey := config.ParseEnv()
	llm, err := config.InitLLM(provider, context.TODO(), apiKey)
	if err != nil {
		log.Fatal(err)
	}

	// connect to LLM and prompt it for an output
	ctx := context.Background()
	if os.Getenv("PROMPT_FIRST") != "" {
		prompt := "What would be a good company name for a company that makes colorful socks? " + rag.LoadGroundingContext()
		completion, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(completion)
	}

	// connect to Weaviate
	log.Println("connecting to Weaviate")
	if err := vectors.InitWeaviate(); err != nil {
		log.Fatal(err)
	}

	log.Println("getting schema")
	_, err = vectors.GetSchema()
	if err != nil {
		log.Fatal(err)
	}

	res, err := vectors.QueryData(ctx, 1)
	if err != nil {
		panic(err)
	}

	if len(res) == 0 {
		if err := vectors.InsertData(ctx); err != nil {
			panic(err)
		}

	}

	generatedData, err := vectors.GenerativeSearch(context.Background(), os.Getenv("QUERY_STRING"), llm)
	if err != nil {
		panic(err)
	}

	log.Println(generatedData)
	log.Println("all done!")

}
