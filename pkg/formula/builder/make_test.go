package builder

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/stream/streams"
)

func TestBuildMake(t *testing.T) {
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	tmpDir := os.TempDir()
	ritHome := filepath.Join(tmpDir, ".rit-make")
	repoPath := filepath.Join(ritHome, "repos", "commons")

	_ = dirManager.Remove(ritHome)
	_ = dirManager.Remove(repoPath)
	_ = dirManager.Create(repoPath)
	zipFile := filepath.Join("..", "..", "..", "testdata", "ritchie-formulas-test.zip")
	_ = streams.Unzip(zipFile, repoPath)

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

			if tt.out.wantErr && got == nil {
				t.Errorf("Run(%s) got %v, want not nil error", tt.name, got)
			}
		})
	}
}
