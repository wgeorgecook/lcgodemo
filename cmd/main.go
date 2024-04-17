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

	// level 1: connect to LLM and prompt it for an output
	ctx := context.Background()
	prompt := "What would be a good company name for a company that makes colorful socks? " + rag.LoadGroundingContext()
	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(completion)

	// connect to Weaviate
	log.Println("connecting to Weaviate")
	if err := vectors.InitWeaviate(); err != nil {
		log.Fatal(err)
	}

	log.Println("getting schema")
	schema, err := vectors.GetSchema()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("schema: %+v", schema)

	data, err := vectors.GetData(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("queriedData: %+v\n", data)

	generatedData, err := vectors.QueryData(context.Background(), os.Getenv("QUERY_STRING"))
	if err != nil {
		panic(err)
	}
	if len(data.Errors) != 0 {
		for _, err := range data.Errors {
			log.Printf("errors in generating data: %v\n", err)
		}
	} else {
		log.Printf("generated data: %+v\n", generatedData)
	}

	log.Println("all done!")

}
