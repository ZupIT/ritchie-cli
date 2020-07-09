package autocomplete

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
)

type repoListerMock struct{}

func (repoListerMock) List() ([]formula.Repository, error) {
	return []formula.Repository{}, nil
}

func TestGenerate(t *testing.T) {
	type in struct {
		shell ShellName
	}

	type out struct {
		err error
	}

	treeMan := tree.NewTreeManager("../../testdata", repoListerMock{}, api.SingleCoreCmds)
	autocomplete := NewGenerator(treeMan)

	tests := []struct {
		name string
		in   *in
		out  *out
	}{
		{
			name: "autocomplete bash",
			in: &in{
				shell: bash,
			},
			out: &out{
				err: nil,
			},
		},
		{
			name: "autocomplete zsh",
			in: &in{
				shell: zsh,
			},
			out: &out{
				err: nil,
			},
		},
		{
			name: "autocomplete fish",
			in: &in{
				shell: fish,
			},
			out: &out{
				err: nil,
			},
		},
		{
			name: "autocomplete powerShell",
			in: &in{
				shell: powerShell,
			},
			out: &out{
				err: nil,
			},
		},
		{
			name: "autocomplete error",
			in: &in{
				shell: "err",
			},
			out: &out{
				err: ErrNotSupported,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := autocomplete.Generate(tt.in.shell, &cobra.Command{})

			if err != tt.out.err {
				t.Errorf("Generator(%s) got %v, want %v", tt.name, err, tt.out.err)
			}

			if tt.out.err == nil && got == "" {
				t.Errorf("Generator(%s) autocomplete is empty", tt.name)
			}
		})
	}
}
