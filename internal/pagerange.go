package internal

import (
	"fmt"
	"strconv"
	"strings"
)

type PageRange struct {
	from int
	to   int
}

func NewPageRange(from int, to int) (*PageRange, error) {
	if from < 1 || to < from {
		return nil, fmt.Errorf("invalid page range: %d-%d", from, to)
	}

	return &PageRange{
		from: from,
		to:   to,
	}, nil
}

func (r *PageRange) From() int {
	return r.from
}

func (r *PageRange) To() int {
	return r.to
}

func (r *PageRange) String() string {
	if r.from == r.to {
		return strconv.Itoa(r.from)
	}

	return fmt.Sprintf("%d-%d", r.from, r.to)
}

func ParseRanges(s string, total int) ([]PageRange, error) {
	var out []PageRange

	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		from, to := 1, total
		if i := strings.Index(part, "-"); i >= 0 {
			low, high := strings.TrimSpace(part[:i]), strings.TrimSpace(part[i+1:])
			var err error

			if low != "" {
				if from, err = strconv.Atoi(low); err != nil {
					return nil, fmt.Errorf("invalid range: %s", part)
				}
			}

			if high != "" {
				if to, err = strconv.Atoi(high); err != nil {
					return nil, fmt.Errorf("invalid range: %s", part)
				}
			}
		} else {
			n, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid page: %s", part)
			}

			from, to = n, n
		}

		if from < 1 || to > total || from > to {
			return nil, fmt.Errorf("range %s is outside 1-%d", part, total)
		}

		out = append(out, PageRange{from: from, to: to})
	}

	if len(out) == 0 {
		return nil, fmt.Errorf("no valid ranges found")
	}

	return out, nil
}

func ChunkRanges(totalPages int, chunkSize int) ([]PageRange, error) {
	if chunkSize < 1 {
		return nil, fmt.Errorf("chunk size must be greater than 0, got %d", chunkSize)
	}

	if totalPages < 1 {
		return nil, fmt.Errorf("pdf has no pages")
	}

	var ranges []PageRange
	for start := 1; start <= totalPages; start += chunkSize {
		end := min(start+chunkSize-1, totalPages)
		ranges = append(ranges, PageRange{from: start, to: end})
	}

	return ranges, nil
}

func UniqueRanges(ranges []PageRange) []PageRange {
	seen := make(map[PageRange]struct{}, len(ranges))
	unique := make([]PageRange, 0, len(ranges))

	for _, r := range ranges {
		if _, exists := seen[r]; !exists {
			seen[r] = struct{}{}
			unique = append(unique, r)
		}
	}

	return unique
}
