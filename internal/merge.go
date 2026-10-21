package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/brottum-brass/pdfmanager/engine"
)

type Merge struct {
	engine     engine.Engine
	InputFiles []string
	OutputFile string
}

func NewMerge(engine engine.Engine, inputFiles []string, outputFile string) *Merge {
	return &Merge{
		engine:     engine,
		InputFiles: inputFiles,
		OutputFile: outputFile,
	}
}

func (m *Merge) All(divider bool) error {
	if len(m.InputFiles) == 0 {
		return fmt.Errorf("need at least two files to merge")
	}

	absoluteOutputFilepath, err := filepath.Abs(m.OutputFile)
	if err != nil {
		return fmt.Errorf("internal error: %w", err)
	}

	for _, f := range m.InputFiles {
		absoluteInputFilepath, err := filepath.Abs(f)

		if err != nil {
			return fmt.Errorf("internal error: %w", err)
		}

		if absoluteInputFilepath == absoluteOutputFilepath {
			return fmt.Errorf("input file %s is the same as output file %s", f, m.OutputFile)
		}
	}

	if err := os.MkdirAll(filepath.Dir(m.OutputFile), 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	return m.engine.Merge(m.InputFiles, m.OutputFile, divider)
}
