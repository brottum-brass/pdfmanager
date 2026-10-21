package engine

// Engine defines the interface for PDF operations.
//
// This is effectively a wrapper around the pdfcpu library, providing methods for counting pages, trimming PDFs, and merging PDFs.
//
// This interface allows for different implementations of PDF engines, enabling flexibility in how PDF operations are performed.
type Engine interface {
	PageCount(path string) (int, error)
	Trim(in string, out string, pages []string) error
	Merge(in []string, out string, divider bool) error
}
