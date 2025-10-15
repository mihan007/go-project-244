// Package code provides functionality for comparing configuration files
// and generating human-readable diffs in various formats.
package code

import (
	"code/parsing"
	"sort"
)

// GenDiff compares two configuration files and returns a string representation
// of the differences. The format parameter controls the output format.
func GenDiff(filepath1, filepath2, format string) (string, error) {
	data1, err := parsing.ParseFile(filepath1)
	if err != nil {
		return "", err
	}
	data2, err := parsing.ParseFile(filepath2)
	if err != nil {
		return "", err
	}

	// Compute the differences
	diff := computeDiff(data1, data2)

	// Get the appropriate formatter
	formatter, err := NewFormatter(format)
	if err != nil {
		return "", err
	}

	// Format and return the result
	return formatter.Format(diff), nil
}

// computeDiff calculates the differences between two data maps
func computeDiff(data1, data2 map[string]interface{}) []DiffEntry {
	// Collect all unique keys
	allKeys := make(map[string]bool)
	for k := range data1 {
		allKeys[k] = true
	}
	for k := range data2 {
		allKeys[k] = true
	}

	// Sort keys for a consistent output
	sortedKeys := make([]string, 0, len(allKeys))
	for k := range allKeys {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	// Build diff entries
	diff := make([]DiffEntry, 0, len(sortedKeys))
	for _, key := range sortedKeys {
		val1, exists1 := data1[key]
		val2, exists2 := data2[key]

		var entry DiffEntry
		entry.Key = key

		switch {
		case !exists1:
			entry.Status = StatusAdded
			entry.NewVal = val2
		case !exists2:
			entry.Status = StatusRemoved
			entry.OldVal = val1
		case val1 != val2:
			entry.Status = StatusChanged
			entry.OldVal = val1
			entry.NewVal = val2
		default:
			entry.Status = StatusUnchanged
			entry.OldVal = val1
		}

		diff = append(diff, entry)
	}

	return diff
}
