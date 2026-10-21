package engine

import (
	"context"
	"fmt"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type PDFCPUEngine struct {
	config *model.Configuration
}

func NewPDFCPUEngine() (*PDFCPUEngine, error) {
	config, err := api.LoadConfiguration(api.ConfigurationOptions{Mode: api.ConfigurationModeStateless})
	if err != nil {
		return nil, fmt.Errorf("loading pdfcpu configuration: %w", err)
	}

	return &PDFCPUEngine{
		config: config,
	}, nil
}

func (e *PDFCPUEngine) PageCount(ctx context.Context, path string) (int, error) {
	return api.PageCountFile(ctx, path)
}

func (e *PDFCPUEngine) Trim(ctx context.Context, inputFile string, outputFile string, pageRanges []string) error {
	return api.TrimFile(ctx, inputFile, outputFile, pageRanges, e.config)
}

func (e *PDFCPUEngine) Merge(ctx context.Context, inputFiles []string, outputFile string, divider bool) error {
	return api.MergeCreateFile(ctx, inputFiles, outputFile, divider, e.config)
}
