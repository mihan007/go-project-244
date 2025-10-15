package code

import (
	"testing"
)

func TestParseFile(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		want     map[string]interface{}
		wantErr  bool
	}{
		{
			name:     "valid json",
			filepath: createTempJSON(t, `{"key": "value", "number": 42}`),
			want:     map[string]interface{}{"key": "value", "number": float64(42)},
		},
		{
			name:     "empty json object",
			filepath: createTempJSON(t, `{}`),
			want:     map[string]interface{}{},
		},
		{
			name:     "nested json",
			filepath: createTempJSON(t, `{"outer": {"inner": "value"}}`),
			want:     map[string]interface{}{"outer": map[string]interface{}{"inner": "value"}},
		},
		{
			name:     "json with array",
			filepath: createTempJSON(t, `{"list": [1, 2, 3]}`),
			want: map[string]interface{}{
				"list": []interface{}{float64(1), float64(2), float64(3)},
			},
		},
		{
			name:     "unsupported extension",
			filepath: createTempFile(t, "file.txt", "content"),
			wantErr:  true,
		},
		{
			name:     "invalid json",
			filepath: createTempFile(t, "file.json", "not valid json"),
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
			got, err := parseFile(tt.filepath)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("parseFile() returned map with %d keys, want %d", len(got), len(tt.want))
					return
				}
				for k, v := range tt.want {
					gotV, ok := got[k]
					if !ok {
						t.Errorf("parseFile() missing key %q", k)
						continue
					}
					if !deepEqual(gotV, v) {
						t.Errorf("parseFile() key %q = %v, want %v", k, gotV, v)
					}
				}
			}
		})
	}
}
