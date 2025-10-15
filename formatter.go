package code

import (
	"fmt"
)

// Formatter defines the interface for different output formats
type Formatter interface {
	Format(diff []DiffEntry) string
}

// DiffEntry represents a single difference between two files
type DiffEntry struct {
	Key    string
	Status DiffStatus
	OldVal interface{}
	NewVal interface{}
}

// DiffStatus represents the type of difference
type DiffStatus int

const (
	StatusUnchanged DiffStatus = iota // Key exists in both with the same value
	StatusAdded                       // Key only exists in file2
	StatusRemoved                     // Key only exists in file1
	StatusChanged                     // Key exists in both with different values
)

// NewFormatter creates a formatter based on the format name
func NewFormatter(format string) (Formatter, error) {
	switch format {
	case "stylish":
		return &FormatterStylish{}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}
