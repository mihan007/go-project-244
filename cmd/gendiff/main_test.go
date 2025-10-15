package main

import (
	"code/helpers"
	"context"
	"strings"
	"testing"

	"github.com/urfave/cli/v3"
)

func TestGenDiff(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		format     string
		wantErr    bool
		wantOutput string
	}{
		{
			name:       "two identical json files with stylish format",
			args:       []string{helpers.CreateTempJSON(t, `{"key": "value"}`), helpers.CreateTempJSON(t, `{"key": "value"}`)},
			format:     "stylish",
			wantErr:    false,
			wantOutput: "{\n    key: value\n}",
		},
		{
			name:       "two different json files",
			args:       []string{helpers.CreateTempJSON(t, `{"key": "old"}`), helpers.CreateTempJSON(t, `{"key": "new"}`)},
			format:     "stylish",
			wantErr:    false,
			wantOutput: "{\n  - key: old\n  + key: new\n}",
		},
		{
			name:    "no arguments provided",
			args:    []string{},
			format:  "stylish",
			wantErr: false,
		},
		{
			name:    "only one argument provided",
			args:    []string{helpers.CreateTempJSON(t, `{}`)},
			format:  "stylish",
			wantErr: true,
		},
		{
			name: "more than two arguments",
			args: []string{
				helpers.CreateTempJSON(t, `{}`),
				helpers.CreateTempJSON(t, `{}`),
				helpers.CreateTempJSON(t, `{}`),
			},
			format:  "stylish",
			wantErr: true,
		},
		{
			name:    "first file does not exist",
			args:    []string{"nonexistent1.json", helpers.CreateTempJSON(t, `{}`)},
			format:  "stylish",
			wantErr: true,
		},
		{
			name:    "second file does not exist",
			args:    []string{helpers.CreateTempJSON(t, `{}`), "nonexistent2.json"},
			format:  "stylish",
			wantErr: true,
		},
		{
			name:    "unsupported format",
			args:    []string{helpers.CreateTempJSON(t, `{}`), helpers.CreateTempJSON(t, `{}`)},
			format:  "plain",
			wantErr: true,
		},
		{
			name: "yaml files",
			args: []string{
				helpers.CreateTempYAML(t, "key: value1"),
				helpers.CreateTempYAML(t, "key: value2"),
			},
			format:     "stylish",
			wantErr:    false,
			wantOutput: "{\n  - key: value1\n  + key: value2\n}",
		},
		{
			name: "mixed json and yaml files",
			args: []string{
				helpers.CreateTempJSON(t, `{"key": "value1"}`),
				helpers.CreateTempYAML(t, "key: value2"),
			},
			format:     "stylish",
			wantErr:    false,
			wantOutput: "{\n  - key: value1\n  + key: value2\n}",
		},
		{
			name: "complex diff with multiple changes",
			args: []string{
				helpers.CreateTempJSON(t, `{"a": 1, "b": 2, "c": 3}`),
				helpers.CreateTempJSON(t, `{"a": 1, "b": 20, "d": 4}`),
			},
			format:     "stylish",
			wantErr:    false,
			wantOutput: "{\n    a: 1\n  - b: 2\n  + b: 20\n  - c: 3\n  + d: 4\n}",
		},
		{
			name:    "invalid json in first file",
			args:    []string{helpers.CreateTempFile(t, "*.json", "invalid json"), helpers.CreateTempJSON(t, `{}`)},
			format:  "stylish",
			wantErr: true,
		},
		{
			name:    "invalid json in second file",
			args:    []string{helpers.CreateTempJSON(t, `{}`), helpers.CreateTempFile(t, "*.json", "invalid json")},
			format:  "stylish",
			wantErr: true,
		},
		{
			name:    "unsupported file extension",
			args:    []string{helpers.CreateTempFile(t, "*.txt", "content"), helpers.CreateTempJSON(t, `{}`)},
			format:  "stylish",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock command with the necessary fields
			cmd := &cli.Command{
				Name: "gendiff",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "format",
						Value: tt.format,
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					return GenDiff(ctx, c)
				},
			}

			// Build full args for the command (program name + actual args)
			fullArgs := append([]string{"gendiff", "--format", tt.format}, tt.args...)

			// Run the command with the arguments
			ctx := context.Background()
			err := cmd.Run(ctx, fullArgs)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenDiff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// For error cases, check the error message
			if tt.wantErr && err != nil && len(tt.args) != 2 && len(tt.args) > 0 {
				expectedErrMsg := "expected 2 arguments"
				if !strings.Contains(err.Error(), expectedErrMsg) {
					t.Errorf("GenDiff() error message = %v, should contain %v", err.Error(), expectedErrMsg)
				}
			}
		})
	}
}

func TestGenDiffWithFlags(t *testing.T) {
	tests := []struct {
		name       string
		file1      string
		file2      string
		formatFlag string
		wantErr    bool
	}{
		{
			name:       "default format (stylish)",
			file1:      helpers.CreateTempJSON(t, `{"key": "value"}`),
			file2:      helpers.CreateTempJSON(t, `{"key": "value"}`),
			formatFlag: "stylish",
			wantErr:    false,
		},
		{
			name:       "explicit stylish format",
			file1:      helpers.CreateTempJSON(t, `{"key": "value"}`),
			file2:      helpers.CreateTempJSON(t, `{"key": "value"}`),
			formatFlag: "stylish",
			wantErr:    false,
		},
		{
			name:       "unsupported format",
			file1:      helpers.CreateTempJSON(t, `{"key": "value"}`),
			file2:      helpers.CreateTempJSON(t, `{"key": "value"}`),
			formatFlag: "json",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cli.Command{
				Name: "gendiff",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "format",
						Value: tt.formatFlag,
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					return GenDiff(ctx, c)
				},
			}

			fullArgs := []string{"gendiff", "--format", tt.formatFlag, tt.file1, tt.file2}
			err := cmd.Run(context.Background(), fullArgs)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenDiff() with format flag error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenDiffEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() (string, string)
		format  string
		wantErr bool
	}{
		{
			name: "empty json files",
			setup: func() (string, string) {
				return helpers.CreateTempJSON(t, `{}`), helpers.CreateTempJSON(t, `{}`)
			},
			format:  "stylish",
			wantErr: false,
		},
		{
			name: "empty yaml files",
			setup: func() (string, string) {
				return helpers.CreateTempYAML(t, "{}"), helpers.CreateTempYAML(t, "{}")
			},
			format:  "stylish",
			wantErr: false,
		},
		{
			name: "json with special characters",
			setup: func() (string, string) {
				return helpers.CreateTempJSON(t, `{"key": "value with spaces"}`),
					helpers.CreateTempJSON(t, `{"key": "value with spaces"}`)
			},
			format:  "stylish",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file1, file2 := tt.setup()

			cmd := &cli.Command{
				Name: "gendiff",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "format",
						Value: tt.format,
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					return GenDiff(ctx, c)
				},
			}

			fullArgs := []string{"gendiff", "--format", tt.format, file1, file2}
			err := cmd.Run(context.Background(), fullArgs)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenDiff() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
