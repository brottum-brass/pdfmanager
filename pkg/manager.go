package pkg

import (
	"fmt"

	"github.com/brottum-brass/pdfmanager/engine"
	"github.com/brottum-brass/pdfmanager/internal"
)

type Manager struct {
	engine        engine.Engine
	outputRootDir string
}

func NewManager(engine engine.Engine, outputRootDir string) (*Manager, error) {
	return &Manager{
		engine:        engine,
		outputRootDir: outputRootDir,
	}, nil
}

func NewDefaultManager() (*Manager, error) {
	pdfcpuEngine, err := engine.NewPDFCPUEngine()
	if err != nil {
		return nil, fmt.Errorf("creating pdfcpu engine: %w", err)
	}

	downloadsDir, err := internal.DefaultDir()
	if err != nil {
		return nil, fmt.Errorf("getting default directory: %w", err)
	}

	return NewManager(pdfcpuEngine, downloadsDir)
}

func (m *Manager) PageCount(path string) (int, error) {
	return m.engine.PageCount(path)
}

func (m *Manager) SplitRanges(inputFile string, outputDirName string, specification string, combine bool) ([]string, error) {
	dir, err := internal.ResolveInputDir(m.outputRootDir, outputDirName)
	if err != nil {
		return nil, err
	}

	return internal.NewSplit(m.engine, inputFile, dir).RangesFromSpecification(specification, combine)
}

func (m *Manager) SplitEveryNPages(inputFile string, outputDirName string, n int) ([]string, error) {
	dir, err := internal.ResolveInputDir(m.outputRootDir, outputDirName)
	if err != nil {
		return nil, err
	}

	return internal.NewSplit(m.engine, inputFile, dir).EveryNPages(n)
}

func (m *Manager) SplitEachPage(inputFile string, outputDirName string) ([]string, error) {
	return m.SplitEveryNPages(inputFile, outputDirName, 1)
}

func (m *Manager) Merge(inputFiles []string, outputFileName string, divider bool) (string, error) {
	if _, err := internal.ResolveInputDir(m.outputRootDir, outputFileName); err != nil {
		return "", fmt.Errorf("output name: %w", err)
	}

	outputFile, err := internal.ResolveInputDir(m.outputRootDir, internal.EnsurePDFExtension(outputFileName))
	if err != nil {
		return "", fmt.Errorf("output name: %w", err)
	}

	merge := internal.NewMerge(m.engine, inputFiles, outputFile)
	if err := merge.All(divider); err != nil {
		return "", err
	}

	return outputFile, nil
}
