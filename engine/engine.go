package engine

type Engine interface {
	PageCount(path string) (int, error)
	Trim(in string, out string, pages []string) error
	Merge(in []string, out string, divider bool) error
}
