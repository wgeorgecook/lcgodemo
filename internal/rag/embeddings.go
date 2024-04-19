package rag

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/openai"
)

const TOKEN_BATCH_SIZE = 500

var (
	embedder embeddings.Embedder
)

func initEmbedder() error {
	log.Println("initializing embedder")
	defer log.Println("initialized!")
	if embedder != nil {
		return nil
	}

	if openAiKey := os.Getenv("OPENAI_API_KEY"); openAiKey == "" {
		return errors.New("no open ai key set")
	}
	llm, err := openai.New(openai.WithEmbeddingModel("text-embedding-ada-002"))
	if err != nil {
		return err
	}
	embedder, err = embeddings.NewEmbedder(llm, embeddings.WithBatchSize(TOKEN_BATCH_SIZE))
	if err != nil {
		return err
	}

	return nil
}

func EmbedChunkedDocument(ctx context.Context, texts []string) ([]float32, error) {
	log.Println("start embed chunked documents")
	if err := initEmbedder(); err != nil {
		return nil, err
	}
	vectors, err := embedder.EmbedDocuments(ctx, texts)
	if err != nil {
		return nil, err
	}
	log.Println("embed done!")
	return vectors[0], nil
}

func EmbedQuery(ctx context.Context, query string) ([]float32, error) {
	if err := initEmbedder(); err != nil {
		return nil, err
	}

	vectors, err := embedder.EmbedQuery(ctx, query)
	if err != nil {
		return nil, err
	}

	return vectors, nil
}

// praise be to Gemini for the basis of this
func CreateStringSliceFromBytes(content []byte) []string {
	// Split the text into lines
	return strings.Split(string(content), "\n")
}
