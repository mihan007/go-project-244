package code

import (
	"code/helpers"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenDiffJSON(t *testing.T) {
	tests := []struct {
		name    string
		file1   string
		file2   string
		format  string
		want    string
		wantErr bool
	}{
		{
			name:   "identical files",
			file1:  helpers.CreateTempJSON(t, `{"key": "value"}`),
			file2:  helpers.CreateTempJSON(t, `{"key": "value"}`),
			format: "stylish",
			want:   "{\n    key: value\n}",
		},
		{
			name:   "added key",
			file1:  helpers.CreateTempJSON(t, `{}`),
			file2:  helpers.CreateTempJSON(t, `{"key": "value"}`),
			format: "stylish",
			want:   "{\n  + key: value\n}",
		},
		{
			name:   "removed key",
			file1:  helpers.CreateTempJSON(t, `{"key": "value"}`),
			file2:  helpers.CreateTempJSON(t, `{}`),
			format: "stylish",
			want:   "{\n  - key: value\n}",
		},
		{
			name:   "changed value",
			file1:  helpers.CreateTempJSON(t, `{"key": "old"}`),
			file2:  helpers.CreateTempJSON(t, `{"key": "new"}`),
			format: "stylish",
			want:   "{\n  - key: old\n  + key: new\n}",
		},
		{
			name:   "multiple changes",
			file1:  helpers.CreateTempJSON(t, `{"a": 1, "b": 2, "c": 3}`),
			file2:  helpers.CreateTempJSON(t, `{"a": 1, "b": 20, "d": 4}`),
			format: "stylish",
			want:   "{\n    a: 1\n  - b: 2\n  + b: 20\n  - c: 3\n  + d: 4\n}",
		},
		{
			name:    "file1 does not exist",
			file1:   "nonexistent.json",
			file2:   helpers.CreateTempJSON(t, `{}`),
			format:  "stylish",
			wantErr: true,
		},
		{
			name:    "file2 does not exist",
			file1:   helpers.CreateTempJSON(t, `{}`),
			file2:   "nonexistent.json",
			format:  "stylish",
			wantErr: true,
		},
		{
			name:    "unsupported format",
			file1:   helpers.CreateTempJSON(t, `{}`),
			file2:   helpers.CreateTempJSON(t, `{}`),
			format:  "plain",
			wantErr: true,
		},
		{
			name:    "unsupported file extension",
			file1:   helpers.CreateTempFile(t, "file.txt", "some content"),
			file2:   helpers.CreateTempJSON(t, `{}`),
			format:  "stylish",
			wantErr: true,
		},
		{
			name:    "invalid json in file1",
			file1:   helpers.CreateTempFile(t, "file.json", "invalid json"),
			file2:   helpers.CreateTempJSON(t, `{}`),
			format:  "stylish",
			wantErr: true,
		},
		{
			name:    "invalid json in file2",
			file1:   helpers.CreateTempJSON(t, `{}`),
			file2:   helpers.CreateTempFile(t, "file.json", "invalid json"),
			format:  "stylish",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenDiff(tt.file1, tt.file2, tt.format)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGenDiffYAML(t *testing.T) {
	tests := []struct {
		name      string
		filepath1 string
		filepath2 string
		format    string
		want      string
		wantErr   bool
	}{
		{
			name:      "both yaml files have same content",
			filepath1: helpers.CreateTempYAML(t, "key: value"),
			filepath2: helpers.CreateTempYAML(t, "key: value"),
			format:    "stylish",
			want:      "{\n    key: value\n}",
		},
		{
			name:      "yaml one key added",
			filepath1: helpers.CreateTempYAML(t, "key1: value1"),
			filepath2: helpers.CreateTempYAML(t, "key1: value1\nkey2: value2"),
			format:    "stylish",
			want:      "{\n    key1: value1\n  + key2: value2\n}",
		},
		{
			name:      "yaml one key removed",
			filepath1: helpers.CreateTempYAML(t, "key1: value1\nkey2: value2"),
			filepath2: helpers.CreateTempYAML(t, "key1: value1"),
			format:    "stylish",
			want:      "{\n    key1: value1\n  - key2: value2\n}",
		},
		{
			name:      "yaml one key changed",
			filepath1: helpers.CreateTempYAML(t, "key: value1"),
			filepath2: helpers.CreateTempYAML(t, "key: value2"),
			format:    "stylish",
			want:      "{\n  - key: value1\n  + key: value2\n}",
		},
		{
			name:      "yaml with numbers",
			filepath1: helpers.CreateTempYAML(t, "timeout: 30\nretries: 3"),
			filepath2: helpers.CreateTempYAML(t, "timeout: 60\nretries: 3"),
			format:    "stylish",
			want:      "{\n    retries: 3\n  - timeout: 30\n  + timeout: 60\n}",
		},
		{
			name:      "yaml with booleans",
			filepath1: helpers.CreateTempYAML(t, "enabled: true\ndebug: false"),
			filepath2: helpers.CreateTempYAML(t, "enabled: false\ndebug: false"),
			format:    "stylish",
			want:      "{\n    debug: false\n  - enabled: true\n  + enabled: false\n}",
		},
		{
			name:      "yaml complex diff",
			filepath1: helpers.CreateTempYAML(t, "host: localhost\nport: 8080\nssl: false"),
			filepath2: helpers.CreateTempYAML(t, "host: localhost\nport: 443\nssl: true\ntimeout: 30"),
			format:    "stylish",
			want:      "{\n    host: localhost\n  - port: 8080\n  + port: 443\n  - ssl: false\n  + ssl: true\n  + timeout: 30\n}",
		},
		{
			name:      "mixed json and yaml files",
			filepath1: helpers.CreateTempJSON(t, `{"key": "value1"}`),
			filepath2: helpers.CreateTempYAML(t, "key: value2"),
			format:    "stylish",
			want:      "{\n  - key: value1\n  + key: value2\n}",
		},
		{
			name:      "yaml file does not exist",
			filepath1: "does-not-exist.yaml",
			filepath2: helpers.CreateTempYAML(t, "key: value"),
			format:    "stylish",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenDiff(tt.filepath1, tt.filepath2, tt.format)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestComputeDiff(t *testing.T) {
	tests := []struct {
		name  string
		data1 map[string]interface{}
		data2 map[string]interface{}
		want  []DiffEntry
	}{
		{
			name:  "empty maps",
			data1: map[string]interface{}{},
			data2: map[string]interface{}{},
			want:  []DiffEntry{},
		},
		{
			name:  "identical maps",
			data1: map[string]interface{}{"key": "value"},
			data2: map[string]interface{}{"key": "value"},
			want: []DiffEntry{
				{Key: "key", Status: StatusUnchanged, OldVal: "value"},
			},
		},
		{
			name:  "added key",
			data1: map[string]interface{}{},
			data2: map[string]interface{}{"key": "value"},
			want: []DiffEntry{
				{Key: "key", Status: StatusAdded, NewVal: "value"},
			},
		},
		{
			name:  "removed key",
			data1: map[string]interface{}{"key": "value"},
			data2: map[string]interface{}{},
			want: []DiffEntry{
				{Key: "key", Status: StatusRemoved, OldVal: "value"},
			},
		},
		{
			name:  "changed value",
			data1: map[string]interface{}{"key": "old"},
			data2: map[string]interface{}{"key": "new"},
			want: []DiffEntry{
				{Key: "key", Status: StatusChanged, OldVal: "old", NewVal: "new"},
			},
		},
		{
			name: "multiple keys with various changes",
			data1: map[string]interface{}{
				"unchanged": "same",
				"removed":   "gone",
				"changed":   "old",
			},
			data2: map[string]interface{}{
				"unchanged": "same",
				"added":     "new",
				"changed":   "new",
			},
			want: []DiffEntry{
				{Key: "added", Status: StatusAdded, NewVal: "new"},
				{Key: "changed", Status: StatusChanged, OldVal: "old", NewVal: "new"},
				{Key: "removed", Status: StatusRemoved, OldVal: "gone"},
				{Key: "unchanged", Status: StatusUnchanged, OldVal: "same"},
			},
		},
		{
			name:  "different types",
			data1: map[string]interface{}{"key": "string"},
			data2: map[string]interface{}{"key": float64(123)},
			want: []DiffEntry{
				{Key: "key", Status: StatusChanged, OldVal: "string", NewVal: float64(123)},
			},
		},
		{
			name:  "boolean values",
			data1: map[string]interface{}{"flag": true},
			data2: map[string]interface{}{"flag": false},
			want: []DiffEntry{
				{Key: "flag", Status: StatusChanged, OldVal: true, NewVal: false},
			},
		},
		{
			name:  "numeric values",
			data1: map[string]interface{}{"count": float64(10)},
			data2: map[string]interface{}{"count": float64(20)},
			want: []DiffEntry{
				{Key: "count", Status: StatusChanged, OldVal: float64(10), NewVal: float64(20)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := computeDiff(tt.data1, tt.data2)
			require.Equal(t, len(tt.want), len(got))
			for i, entry := range got {
				assert.Equal(t, tt.want[i].Key, entry.Key)
				assert.Equal(t, tt.want[i].Status, entry.Status)
				assert.Equal(t, tt.want[i].OldVal, entry.OldVal)
				assert.Equal(t, tt.want[i].NewVal, entry.NewVal)
			}
		})
	}
}
