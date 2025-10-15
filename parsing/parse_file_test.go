package parsing

import (
	"code/helpers"
	"testing"
)

func TestParseJSON(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		want     map[string]interface{}
		wantErr  bool
	}{
		{
			name:     "valid json",
			filepath: helpers.CreateTempJSON(t, `{"key": "value", "number": 42}`),
			want:     map[string]interface{}{"key": "value", "number": float64(42)},
		},
		{
			name:     "empty json object",
			filepath: helpers.CreateTempJSON(t, `{}`),
			want:     map[string]interface{}{},
		},
		{
			name:     "nested json",
			filepath: helpers.CreateTempJSON(t, `{"outer": {"inner": "value"}}`),
			want:     map[string]interface{}{"outer": map[string]interface{}{"inner": "value"}},
		},
		{
			name:     "json with array",
			filepath: helpers.CreateTempJSON(t, `{"list": [1, 2, 3]}`),
			want: map[string]interface{}{
				"list": []interface{}{float64(1), float64(2), float64(3)},
			},
		},
		{
			name:     "unsupported extension",
			filepath: helpers.CreateTempFile(t, "file.txt", "content"),
			wantErr:  true,
		},
		{
			name:     "invalid json",
			filepath: helpers.CreateTempFile(t, "file.json", "not valid json"),
			wantErr:  true,
		},
		{
			name:     "non-existent file",
			filepath: "does-not-exist.json",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFile(tt.filepath)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("parseJSON() returned map with %d keys, want %d", len(got), len(tt.want))
					return
				}
				for k, v := range tt.want {
					gotV, ok := got[k]
					if !ok {
						t.Errorf("parseJSON() missing key %q", k)
						continue
					}
					if !helpers.DeepEqual(gotV, v) {
						t.Errorf("parseJSON() key %q = %v, want %v", k, gotV, v)
					}
				}
			}
		})
	}
}

func TestParseYAML(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		want     map[string]interface{}
		wantErr  bool
	}{
		{
			name:     "valid yaml",
			filepath: helpers.CreateTempYAML(t, "key: value\nnumber: 42"),
			want:     map[string]interface{}{"key": "value", "number": 42},
		},
		{
			name:     "empty yaml",
			filepath: helpers.CreateTempYAML(t, "{}"),
			want:     map[string]interface{}{},
		},
		{
			name:     "nested yaml",
			filepath: helpers.CreateTempYAML(t, "outer:\n  inner: value"),
			want:     map[string]interface{}{"outer": map[string]interface{}{"inner": "value"}},
		},
		{
			name:     "yaml with array",
			filepath: helpers.CreateTempYAML(t, "list:\n  - 1\n  - 2\n  - 3"),
			want: map[string]interface{}{
				"list": []interface{}{1, 2, 3},
			},
		},
		{
			name:     "yaml with boolean",
			filepath: helpers.CreateTempYAML(t, "enabled: true\ndisabled: false"),
			want:     map[string]interface{}{"enabled": true, "disabled": false},
		},
		{
			name:     "yaml with null",
			filepath: helpers.CreateTempYAML(t, "key: null"),
			want:     map[string]interface{}{"key": nil},
		},
		{
			name:     "invalid yaml",
			filepath: helpers.CreateTempFile(t, "file.yaml", "key: value\n  invalid: indentation"),
			wantErr:  true,
		},
		{
			name:     "non-existent file",
			filepath: "does-not-exist.yaml",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFile(tt.filepath)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("parseYAML() returned map with %d keys, want %d", len(got), len(tt.want))
					return
				}
				for k, v := range tt.want {
					gotV, ok := got[k]
					if !ok {
						t.Errorf("parseYAML() missing key %q", k)
						continue
					}
					if !helpers.DeepEqual(gotV, v) {
						t.Errorf("parseYAML() key %q = %v, want %v", k, gotV, v)
					}
				}
			}
		})
	}
}
