package builder

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildMake(t *testing.T) {
	tmpDir := os.TempDir()
	ritHome := filepath.Join(tmpDir, ".rit-builder")
	repoPath := filepath.Join(ritHome, "repos", "commons")

	buildMake := NewBuildMake()

	type in struct {
		formPath string
	}

	type out struct {
		wantErr bool
	}

	tests := []struct {
		name string
		in   in
		out  out
	}{
		{
			name: "success",
			in: in{
				formPath: filepath.Join(repoPath, "testing", "formula"),
			},
			out: out{wantErr: false},
		},
		{
			name: "makefile error",
			in: in{
				formPath: repoPath,
			},
			out: out{wantErr: true},
		},
		{
			name: "Chdir error",
			in: in{
				formPath: filepath.Join(repoPath, "invalid"),
			},
			out: out{wantErr: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildMake.Build(tt.in.formPath)

			if got != nil && !tt.out.wantErr {
				t.Errorf("Run(%s) got %v, want not nil error", tt.name, got)
			}
		})
	}
}
