package pdf

import (
	"context"
	"log"

	"github.com/rudolfoborges/pdf2go"
)

func ReadFromFile(ctx context.Context) ([]string, error) {
	pdf, err := pdf2go.New("imports/file.pdf", pdf2go.Config{
		LogLevel: pdf2go.LogLevelError,
	})

	if err != nil {
		return nil, err
	}

	pages, err := pdf.Pages()

	if err != nil {
		return nil, err
	}

	var pageArr []string
	for _, page := range pages {
		text, err := page.Text()
		if err != nil {
			log.Printf("could not read page text: %v\n", err)
			continue
		}
		pageArr = append(pageArr, text)
	}

	return pageArr, nil
}
