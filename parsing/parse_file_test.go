package parsing

import (
	"code/helpers"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, len(tt.want), len(got))
			for k, v := range tt.want {
				assert.Contains(t, got, k)
				assert.True(t, helpers.DeepEqual(got[k], v))
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
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, len(tt.want), len(got))
			for k, v := range tt.want {
				assert.Contains(t, got, k)
				assert.True(t, helpers.DeepEqual(got[k], v))
			}
		})
	}
}
