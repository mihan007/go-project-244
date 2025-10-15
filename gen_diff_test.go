package code

import (
	"testing"
)

func TestGenDiff(t *testing.T) {
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
			file1:  createTempJSON(t, `{"key": "value"}`),
			file2:  createTempJSON(t, `{"key": "value"}`),
			format: "stylish",
			want:   "{\n    key: value\n}",
		},
		{
			name:   "added key",
			file1:  createTempJSON(t, `{}`),
			file2:  createTempJSON(t, `{"key": "value"}`),
			format: "stylish",
			want:   "{\n  + key: value\n}",
		},
		{
			name:   "removed key",
			file1:  createTempJSON(t, `{"key": "value"}`),
			file2:  createTempJSON(t, `{}`),
			format: "stylish",
			want:   "{\n  - key: value\n}",
		},
		{
			name:   "changed value",
			file1:  createTempJSON(t, `{"key": "old"}`),
			file2:  createTempJSON(t, `{"key": "new"}`),
			format: "stylish",
			want:   "{\n  - key: old\n  + key: new\n}",
		},
		{
			name:   "multiple changes",
			file1:  createTempJSON(t, `{"a": 1, "b": 2, "c": 3}`),
			file2:  createTempJSON(t, `{"a": 1, "b": 20, "d": 4}`),
			format: "stylish",
			want:   "{\n    a: 1\n  - b: 2\n  + b: 20\n  - c: 3\n  + d: 4\n}",
		},
		{
			name:    "file1 does not exist",
			file1:   "nonexistent.json",
			file2:   createTempJSON(t, `{}`),
			format:  "stylish",
			wantErr: true,
		},
		{
			name:    "file2 does not exist",
			file1:   createTempJSON(t, `{}`),
			file2:   "nonexistent.json",
			format:  "stylish",
			wantErr: true,
		},
		{
			name:    "unsupported format",
			file1:   createTempJSON(t, `{}`),
			file2:   createTempJSON(t, `{}`),
			format:  "plain",
			wantErr: true,
		},
		{
			name:    "unsupported file extension",
			file1:   createTempFile(t, "file.txt", "some content"),
			file2:   createTempJSON(t, `{}`),
			format:  "stylish",
			wantErr: true,
		},
		{
			name:    "invalid json in file1",
			file1:   createTempFile(t, "file.json", "invalid json"),
			file2:   createTempJSON(t, `{}`),
			format:  "stylish",
			wantErr: true,
		},
		{
			name:    "invalid json in file2",
			file1:   createTempJSON(t, `{}`),
			file2:   createTempFile(t, "file.json", "invalid json"),
			format:  "stylish",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenDiff(tt.file1, tt.file2, tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenDiff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("GenDiff() =\n%v\nwant:\n%v", got, tt.want)
			}
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
			if len(got) != len(tt.want) {
				t.Errorf("computeDiff() returned %d entries, want %d", len(got), len(tt.want))
				return
			}
			for i, entry := range got {
				if entry.Key != tt.want[i].Key {
					t.Errorf("Entry %d: Key = %v, want %v", i, entry.Key, tt.want[i].Key)
				}
				if entry.Status != tt.want[i].Status {
					t.Errorf("Entry %d: Status = %v, want %v", i, entry.Status, tt.want[i].Status)
				}
				if entry.OldVal != tt.want[i].OldVal {
					t.Errorf("Entry %d: OldVal = %v, want %v", i, entry.OldVal, tt.want[i].OldVal)
				}
				if entry.NewVal != tt.want[i].NewVal {
					t.Errorf("Entry %d: NewVal = %v, want %v", i, entry.NewVal, tt.want[i].NewVal)
				}
			}
		})
	}
}
