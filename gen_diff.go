package code

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sort"
)

func GenDiff(filepath1, filepath2, format string) (string, error) {
	parsedJson1, err := parseFile(filepath1)
	if err != nil {
		return "", err
	}
	parsedJson2, err := parseFile(filepath2)
	if err != nil {
		return "", err
	}
	allKeys := make(map[string]bool)
	for k := range parsedJson1 {
		allKeys[k] = true
	}
	for k := range parsedJson2 {
		allKeys[k] = true
	}
	sortedKeys := make([]string, 0, len(allKeys))
	for k := range allKeys {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	result := "{\n"
	for _, key := range sortedKeys {
		val1, exists1 := parsedJson1[key]
		val2, exists2 := parsedJson2[key]

		switch {
		case !exists1:
			result += fmt.Sprintf("  + %s: %v\n", key, val2)
		case !exists2:
			result += fmt.Sprintf("  - %s: %v\n", key, val1)
		case val1 != val2:
			result += fmt.Sprintf("  - %s: %v\n", key, val1)
			result += fmt.Sprintf("  + %s: %v\n", key, val2)
		default:
			result += fmt.Sprintf("    %s: %v\n", key, val1)
		}
	}
	result += "}"
	return result, nil
}

func parseFile(filepath string) (map[string]interface{}, error) {
	ext := path.Ext(filepath)
	if ext == ".json" {
		jsonData, err := os.ReadFile(filepath)
		if err != nil {
			return nil, err
		}
		return ParseJson(jsonData)
	}
	return nil, errors.New("unknown format")
}
