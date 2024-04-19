package vectors

import (
	"context"
	"errors"
	"fmt"
	"lcgodemo/internal/pdf"
	"lcgodemo/internal/rag"
	"log"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/schema"
	"github.com/weaviate/weaviate/entities/models"
)

var client *weaviate.Client

func GetSchema() (*schema.Dump, error) {
	if err := CreateSchemaIfNotExists(); err != nil {
		return nil, err
	}

	schema, err := client.Schema().Getter().Do(context.Background())
	if err != nil {
		return nil, err
	}

	return schema, nil
}

func InitWeaviate() error {
	if client != nil {
		return nil
	}

	cfg := weaviate.Config{
		Host:   "weaviate:8080",
		Scheme: "http",
		Headers: map[string]string{
			"X-OpenAI-Api-Key": os.Getenv("OPENAI_API_KEY"),
		},
	}

	var err error
	client, err = weaviate.NewClient(cfg)
	if err != nil {
		return err
	}

	return nil
}

func InsertData(ctx context.Context) error {
	log.Println("inserting data")
	defer log.Println("done!")
	content, err := pdf.ReadFromFile(ctx)
	if err != nil {
		return err
	}
	vectors, err := rag.EmbedChunkedDocument(ctx, content)
	if err != nil {
		return err
	}

	data := &models.Object{
		Class: pdfClass,
		Properties: map[string]any{
			"title":   "file.pdf",
			"content": content,
		},
		Vector: vectors,
	}

	log.Println("start batch insert")
	defer log.Println("batch done!")
	batchRes, err := client.Batch().ObjectsBatcher().WithObjects(data).Do(ctx)
	if err != nil {
		return err
	}

	var errors []string
	for _, res := range batchRes {
		if res.Result.Errors != nil {
			for _, err := range res.Result.Errors.Error {
				errors = append(errors, fmt.Sprintf("%v, ", err.Message))
			}

		}
	}

	if len(errors) != 0 {
		return fmt.Errorf("error in insert: %v", errors)
	}
	return nil
}

func QueryData(ctx context.Context, limit int) ([]*models.Object, error) {

	result, err := client.Data().ObjectsGetter().
		WithClassName(pdfClass).
		WithLimit(1).
		WithVector().
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GenerativeSearch(ctx context.Context, search string, llm llms.Model) ([]string, error) {
	queryVector, err := rag.EmbedQuery(ctx, search)
	if err != nil {
		return nil, err
	}

	nearVector := client.GraphQL().
		NearVectorArgBuilder().
		WithVector(queryVector)

	fields := []graphql.Field{
		{
			Name: "content",
		},
	}

	results, err := client.GraphQL().Get().
		WithClassName(pdfClass).
		WithFields(fields...).
		WithLimit(1).
		WithNearVector(nearVector).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	if len(results.Errors) > 0 {
		return nil, errors.New(results.Errors[0].Message)
	}

	var generatedPrompts []string
	for _, res := range results.Data {
		var prompt = fmt.Sprintf("given the following, summarize the content in a 256 character tweet: %s", res)
		result, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
		if err != nil {
			log.Printf("cannot generate prompt: %s", err.Error())
		}
		generatedPrompts = append(generatedPrompts, result)
	}

	return generatedPrompts, nil
}
