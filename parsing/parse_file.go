package parsing

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

func ParseFile(filepath string) (map[string]interface{}, error) {
	ext := path.Ext(filepath)
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	switch ext {
	case ".json":
		return parseJSON(data)
	case ".yaml", ".yml":
		return parseYAML(data)
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

func parseYAML(yamlData []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := yaml.Unmarshal(yamlData, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse yaml: %w", err)
	}

	return result, nil
}
