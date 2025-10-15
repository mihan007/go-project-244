package code

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

func parseFile(filepath string) (map[string]interface{}, error) {
	ext := path.Ext(filepath)
	if ext == ".json" {
		jsonData, err := os.ReadFile(filepath)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}
		return parseJSON(jsonData)
	}
	return nil, fmt.Errorf("unsupported file extension: %s", ext)
}

func parseJSON(jsonData []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	return result, nil
}
