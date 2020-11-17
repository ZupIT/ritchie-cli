package builder

import (
	"path/filepath"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

func TestBuildMake(t *testing.T) {
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
			info := formula.BuildInfo{FormulaPath: tt.in.formPath}
			got := buildMake.Build(info)

			if got != nil && !tt.out.wantErr {
				t.Errorf("Run(%s) got %v, want not nil error", tt.name, got)
			}
		})
	}
}
