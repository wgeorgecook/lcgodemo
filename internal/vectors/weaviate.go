package vectors

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"lcgodemo/internal/pdf"
	"lcgodemo/internal/rag"
	"log"

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
	if err := InsertData(context.Background()); err != nil {
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
		// Headers: map[string]string{
		// 	"X-OpenAI-Api-Key": os.Getenv("OPENAI_API_KEY"),
		// },
	}

	var err error
	client, err = weaviate.NewClient(cfg)
	if err != nil {
		return err
	}

	return nil
}

func InsertData(ctx context.Context) error {
	content, err := pdf.ReadFromFile(ctx)
	if err != nil {
		return err
	}
	vectors, err := rag.EmbedChunkedDocument(ctx, content)
	if err != nil {
		return err
	}

	vectorMap := make(map[string]models.Vector)
	for i, text := range vectors {
		page := fmt.Sprintf("page %v", i+1)
		vectorMap[page] = text
	}

	data := &models.Object{
		Class: pdfClass,
		Properties: map[string]any{
			"title":   "file.pdf",
			"content": content,
		},
		Vectors: vectorMap,
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

func GetData(ctx context.Context) (*models.GraphQLResponse, error) {
	fields := []graphql.Field{
		{
			Name: "title",
		},
		{
			Name: "content",
		},
	}
	res, err := client.GraphQL().Get().WithClassName(pdfClass).WithFields(fields...).Do(ctx)
	if err != nil {
		return nil, err
	}

	if len(res.Errors) != 0 {
		for _, err := range res.Errors {
			log.Printf("query error: %+v", err)
		}

		return nil, errors.New("too many errors")
	}

	return res, nil
}
func QueryData(ctx context.Context, search string) (*models.GraphQLResponse, error) {
	fields := []graphql.Field{
		{
			Name: "title",
		},
		{
			Name: "content",
		},
	}

	queryVector, err := rag.EmbedQuery(ctx, search)
	if err != nil {
		return nil, err
	}

	nearText := client.GraphQL().
		NearVectorArgBuilder().
		WithVector(queryVector)

	generativeSearch := graphql.
		NewGenerativeSearch().
		SingleResult("Explain {answer} as you might to a five-year-old.")

	result, err := client.GraphQL().Get().
		WithClassName(pdfClass).
		WithFields(fields...).
		WithNearVector(nearText).
		WithLimit(2).
		WithGenerativeSearch(generativeSearch).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	jsonOutput, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, err
	}
	log.Println(string(jsonOutput))

	return result, nil
}
