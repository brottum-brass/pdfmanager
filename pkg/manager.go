package pkg

import (
	"context"
	"fmt"

	"github.com/brottum-brass/pdfmanager/engine"
	"github.com/brottum-brass/pdfmanager/internal"
)

type Manager struct {
	engine        engine.Engine
	outputRootDir string
}

// NewManager creates a new Manager with the given engine and output root directory.
//
// This function takes two parameters:
//
// - engine: An instance of the Engine interface, which provides methods for PDF operations.
//
// - outputRootDir: A string representing the root directory where output files will be stored.
//
// It returns a pointer to a Manager instance and an error if any issues occur during initialization.
func NewManager(engine engine.Engine, outputRootDir string) (*Manager, error) {
	return &Manager{
		engine:        engine,
		outputRootDir: outputRootDir,
	}, nil
}

// NewDefaultManager creates a new Manager with the default PDF engine and output directory.
//
// This function initializes a PDFCPUEngine and retrieves the default output directory.
//
// It returns a pointer to a Manager instance and an error if any issues occur during initialization.
func NewDefaultManager() (*Manager, error) {
	pdfcpuEngine, err := engine.NewPDFCPUEngine()
	if err != nil {
		return nil, fmt.Errorf("creating pdfcpu engine: %w", err)
	}

	defaultDir, err := internal.DefaultDir()
	if err != nil {
		return nil, fmt.Errorf("getting default directory: %w", err)
	}

	return NewManager(pdfcpuEngine, defaultDir)
}

// PageCount returns the number of pages in the specified PDF file.
//
// This function takes a single parameter:
//
// - ctx: A context for the operation.
//
// - path: A string representing the path to the PDF file.
//
// It returns the number of pages in the PDF file and an error if any issues occur during the operation.
func (m *Manager) PageCount(ctx context.Context, path string) (int, error) {
	return m.engine.PageCount(ctx, path)
}

// SplitRanges splits a PDF file into multiple files based on the provided ranges.
//
// This function takes four parameters:
//
// - ctx: A context for the operation.
//
// - inputFile: A string representing the path to the input PDF file.
//
// - outputDirName: A string representing the name of the output directory.
//
// - specification: A string representing the ranges to split the PDF file by.
//
// - combine: A boolean indicating whether to combine the split files.
//
// It returns a slice of strings representing the paths to the output files and an error if any issues occur during the operation.
func (m *Manager) SplitRanges(ctx context.Context, inputFile string, outputDirName string, specification string, combine bool) ([]string, error) {
	dir, err := internal.ResolveInputDir(m.outputRootDir, outputDirName)
	if err != nil {
		return nil, err
	}

	return internal.NewSplit(m.engine, inputFile, dir).RangesFromSpecification(ctx, specification, combine)
}

// SplitEveryNPages splits a PDF file into multiple files, each containing N pages.
//
// This function takes three parameters:
//
// - ctx: A context for the operation.
//
// - inputFile: A string representing the path to the input PDF file.
//
// - outputDirName: A string representing the name of the output directory.
//
// - n: An integer representing the number of pages in each split file.
//
// It returns a slice of strings representing the paths to the output files and an error if any issues occur during the operation.
func (m *Manager) SplitEveryNPages(ctx context.Context, inputFile string, outputDirName string, n int) ([]string, error) {
	dir, err := internal.ResolveInputDir(m.outputRootDir, outputDirName)
	if err != nil {
		return nil, err
	}

	return internal.NewSplit(m.engine, inputFile, dir).EveryNPages(ctx, n)
}

// SplitEachPage splits a PDF file into multiple files, each containing one page.
//
// This function takes two parameters:
//
// - ctx: A context for the operation.
//
// - inputFile: A string representing the path to the input PDF file.
//
// - outputDirName: A string representing the name of the output directory.
//
// It returns a slice of strings representing the paths to the output files and an error if any issues occur during the operation.
func (m *Manager) SplitEachPage(ctx context.Context, inputFile string, outputDirName string) ([]string, error) {
	return m.SplitEveryNPages(ctx, inputFile, outputDirName, 1)
}

// Merge merges multiple PDF files into a single output file.
//
// This function takes three parameters:
//
// - ctx: A context for the operation.
//
// - inputFiles: A slice of strings representing the paths to the input PDF files.
//
// - outputFileName: A string representing the name of the output PDF file.
//
// - divider: A boolean indicating whether to insert a divider page between merged files.
//
// It returns a string representing the path to the output file and an error if any issues occur during the operation.
func (m *Manager) Merge(ctx context.Context, inputFiles []string, outputFileName string, divider bool) (string, error) {
	if _, err := internal.ResolveInputDir(m.outputRootDir, outputFileName); err != nil {
		return "", fmt.Errorf("output name: %w", err)
	}

	outputFile, err := internal.ResolveInputDir(m.outputRootDir, internal.EnsurePDFExtension(outputFileName))
	if err != nil {
		return "", fmt.Errorf("output name: %w", err)
	}

	merge := internal.NewMerge(m.engine, inputFiles, outputFile)
	if err := merge.All(ctx, divider); err != nil {
		return "", err
	}

	return outputFile, nil
}
