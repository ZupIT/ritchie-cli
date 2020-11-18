package builder

import (
	"path/filepath"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestDockerBuild(t *testing.T) {
	const dockerImg = "ritclizup/rit-go-runner"
	fileManager := stream.NewFileManager()
	buildDocker := NewBuildDocker(fileManager)

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
			name: "docker error",
			in: in{
				formPath: repoPath,
			},
			out: out{wantErr: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := formula.BuildInfo{FormulaPath: tt.in.formPath, DockerImg: dockerImg}
			got := buildDocker.Build(info)

			if tt.out.wantErr && got == nil {
				t.Errorf("Run(%s) got %v, want not nil error", tt.name, got)
			}
		})
	}

}
