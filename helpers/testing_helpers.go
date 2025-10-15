package helpers

import (
	"os"
	"testing"
)

// CreateTempJSON creates a temporary JSON file with the given content
func CreateTempJSON(t *testing.T, content string) string {
	return CreateTempFile(t, "*.json", content)
}

// CreateTempYAML creates a temporary YAML file with the given content
func CreateTempYAML(t *testing.T, content string) string {
	return CreateTempFile(t, "*.yaml", content)
}

// CreateTempFile creates a temporary file with the given pattern and content
func CreateTempFile(t *testing.T, pattern string, content string) string {
	t.Helper()
	tmpFile, err := os.CreateTemp("", pattern)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() {
		if err := tmpFile.Close(); err != nil {
			t.Errorf("Failed to close temp file: %v", err)
		}
	}()

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	t.Cleanup(func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			t.Errorf("Failed to remove temp file: %v", err)
		}
	})

	return tmpFile.Name()
}

// DeepEqual compares two interface{} values, handling nested maps and slices
func DeepEqual(a, b interface{}) bool {
	switch aVal := a.(type) {
	case map[string]interface{}:
		bVal, ok := b.(map[string]interface{})
		if !ok || len(aVal) != len(bVal) {
			return false
		}
		for k, v := range aVal {
			if !DeepEqual(v, bVal[k]) {
				return false
			}
		}
		return true
	case []interface{}:
		bVal, ok := b.([]interface{})
		if !ok || len(aVal) != len(bVal) {
			return false
		}
		for i, v := range aVal {
			if !DeepEqual(v, bVal[i]) {
				return false
			}
		}
		return true
	default:
		return a == b
	}
}
