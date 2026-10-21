package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/brottum-brass/pdfmanager/engine"
)

type Split struct {
	engine        engine.Engine
	InputFile     string
	OutputDirName string
}

func NewSplit(engine engine.Engine, inputFile string, outputDirName string) *Split {
	return &Split{
		engine:        engine,
		InputFile:     inputFile,
		OutputDirName: outputDirName,
	}
}

func (s *Split) baseName() string {
	return strings.TrimSuffix(filepath.Base(s.InputFile), filepath.Ext(s.InputFile))
}

func (s *Split) Ranges(ranges []PageRange, combine bool) ([]string, error) {
	if len(ranges) == 0 {
		return nil, fmt.Errorf("no page ranges provided")
	}

	if err := os.MkdirAll(s.OutputDirName, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	base := s.baseName()

	if combine {
		parts := make([]string, len(ranges))

		for i, r := range ranges {
			parts[i] = r.String()
		}

		outputFile := filepath.Join(s.OutputDirName, fmt.Sprintf("%s_selected.pdf", base))
		if err := s.engine.Trim(s.InputFile, outputFile, parts); err != nil {
			return nil, fmt.Errorf("writing %s: %w", outputFile, err)
		}

		return []string{outputFile}, nil
	}

	uniqueRanges := UniqueRanges(ranges)
	writtenFiles := make([]string, 0, len(uniqueRanges))

	for _, r := range uniqueRanges {
		outputFile := filepath.Join(s.OutputDirName, fmt.Sprintf("%s_%s.pdf", base, r.String()))

		if err := s.engine.Trim(s.InputFile, outputFile, []string{r.String()}); err != nil {
			return nil, fmt.Errorf("writing %s: %w", outputFile, err)
		}

		writtenFiles = append(writtenFiles, outputFile)
	}

	return writtenFiles, nil
}

func (s *Split) RangesFromSpecification(specification string, combine bool) ([]string, error) {
	totalPages, err := s.engine.PageCount(s.InputFile)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", s.InputFile, err)
	}

	ranges, err := ParseRanges(specification, totalPages)
	if err != nil {
		return nil, fmt.Errorf("parsing specification ranges: %w", err)
	}

	return s.Ranges(ranges, combine)
}

func (s *Split) EveryNPages(n int) ([]string, error) {
	totalPages, err := s.engine.PageCount(s.InputFile)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", s.InputFile, err)
	}

	ranges, err := ChunkRanges(totalPages, n)
	if err != nil {
		return nil, fmt.Errorf("chunking pages: %w", err)
	}

	return s.Ranges(ranges, false)
}
