package engine

import "context"

// Engine defines the interface for PDF operations.
//
// This is effectively a wrapper around the pdfcpu library, providing methods for counting pages, trimming PDFs, and merging PDFs.
//
// This interface allows for different implementations of PDF engines, enabling flexibility in how PDF operations are performed.
type Engine interface {
	PageCount(ctx context.Context, path string) (int, error)
	Trim(ctx context.Context, in string, out string, pages []string) error
	Merge(ctx context.Context, in []string, out string, divider bool) error
}
