package code

import (
	"fmt"
	"strings"
)

// FormatterStylish implements the stylish format
type FormatterStylish struct{}

func (f *FormatterStylish) Format(diff []DiffEntry) string {
	var result strings.Builder
	result.WriteString("{\n")

	for _, entry := range diff {
		switch entry.Status {
		case StatusAdded:
			result.WriteString(fmt.Sprintf("  + %s: %v\n", entry.Key, entry.NewVal))
		case StatusRemoved:
			result.WriteString(fmt.Sprintf("  - %s: %v\n", entry.Key, entry.OldVal))
		case StatusChanged:
			result.WriteString(fmt.Sprintf("  - %s: %v\n", entry.Key, entry.OldVal))
			result.WriteString(fmt.Sprintf("  + %s: %v\n", entry.Key, entry.NewVal))
		case StatusUnchanged:
			result.WriteString(fmt.Sprintf("    %s: %v\n", entry.Key, entry.OldVal))
		}
	}

	result.WriteString("}")
	return result.String()
}
