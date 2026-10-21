package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	DEFAULT_OUTPUT_DIR = "Downloads"
	PDF_EXTENSION      = ".pdf"
)

func HomeDir() (string, error) {
	return os.UserHomeDir()
}

func DefaultDir() (string, error) {
	home, err := HomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(home, DEFAULT_OUTPUT_DIR), nil
}

func ResolveInputDir(baseDir string, name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" || name == "." || name == ".." || strings.ContainsAny(name, `/\`) {
		return "", fmt.Errorf("-in must be a plain name without slashes, got %q", name)
	}

	return filepath.Join(baseDir, name), nil
}

func EnsurePDFExtension(name string) string {
	name = strings.TrimSpace(name)
	if strings.EqualFold(filepath.Ext(name), PDF_EXTENSION) {
		return name
	}

	return fmt.Sprintf("%s%s", name, PDF_EXTENSION)
}
