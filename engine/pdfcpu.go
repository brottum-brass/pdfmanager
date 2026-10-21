package engine

import (
	"context"
	"fmt"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type PDFCPUEngine struct {
	config  *model.Configuration
	context context.Context
}

func NewPDFCPUEngine() (*PDFCPUEngine, error) {
	config, err := api.LoadConfiguration(api.ConfigurationOptions{Mode: api.ConfigurationModeStateless})
	if err != nil {
		return nil, fmt.Errorf("loading pdfcpu configuration: %w", err)
	}

	return &PDFCPUEngine{
		config:  config,
		context: context.Background(),
	}, nil
}

func (e *PDFCPUEngine) PageCount(path string) (int, error) {
	return api.PageCountFile(e.context, path)
}

func (e *PDFCPUEngine) Trim(inputFile string, outputFile string, pageRanges []string) error {
	return api.TrimFile(e.context, inputFile, outputFile, pageRanges, e.config)
}

func (e *PDFCPUEngine) Merge(inputFiles []string, outputFile string, divider bool) error {
	return api.MergeCreateFile(e.context, inputFiles, outputFile, divider, e.config)
}
