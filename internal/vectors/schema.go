package vectors

import (
	"context"
	"log"

	"github.com/weaviate/weaviate/entities/models"
)

const (
	pdfClass = "PDF"
)

var PDFSchema = models.Class{
	Class:       pdfClass,
	Description: "Schema for holding vectorized PDF data",
	Properties: []*models.Property{
		{
			Name:        "content",
			Description: "string content of the data",
			DataType:    []string{"text[]"},
		},
		{
			Name:        "title",
			Description: "title of the provided data file",
			DataType:    []string{"text"},
		},
	},
}

func CreateSchemaIfNotExists() error {
	ok, err := client.Schema().ClassExistenceChecker().WithClassName(pdfClass).Do(context.Background())
	if err != nil {
		log.Printf("could not check for class existence: %v\n", err)
	}

	if ok {
		log.Println("class exists, exiting")
		return nil
	}
	log.Println("class does not exist, creating")
	creator := client.Schema().ClassCreator().WithClass(&PDFSchema)
	if err := creator.Do(context.Background()); err != nil {
		return err
	}
	log.Println("created")

	return nil
}
