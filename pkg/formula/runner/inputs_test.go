package runner

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

func TestInputManager_Inputs(t *testing.T) {
	//Todo fix
	t.Skip()
	def := formula.Definition{
		Path: "mock/test",
	}

	home := os.TempDir()
	_ = fileutil.RemoveDir(home + "/formulas")
	defaultSetup := NewDefaultSetup(home)
	preRunner := NewDefaultPreRunner(defaultSetup)
	setup, err := preRunner.PreRun(def)
	if err != nil {
		t.Fatal(err)
	}

	type in struct {
		iText  inputMock
		iList  inputMock
		iBool  inputMock
		iPass  inputMock
		inType api.TermInputType
		stdin  string
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "success stdin",
			in: in{
				iText:  inputMock{text: formula.DefaultCacheNewLabel},
				iList:  inputMock{text: "test"},
				iBool:  inputMock{boolean: false},
				iPass:  inputMock{text: "******"},
				inType: api.Stdin,
				stdin:  `{"sample_text":"test_text","sample_list":"test_list","sample_bool": false}`,
			},
			want: nil,
		},
		{
			name: "success prompt",
			in: in{
				iText:  inputMock{text: formula.DefaultCacheNewLabel},
				iList:  inputMock{text: "test"},
				iBool:  inputMock{boolean: false},
				iPass:  inputMock{text: "******"},
				inType: api.Prompt,
			},
			want: nil,
		},
		{
			name: "error unknown prompt",
			in: in{
				iText:  inputMock{text: formula.DefaultCacheNewLabel},
				iList:  inputMock{text: "test"},
				iBool:  inputMock{boolean: false},
				iPass:  inputMock{text: "******"},
				inType: api.TermInputType(3),
			},
			want: ErrInputNotRecognized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolvers := env.Resolvers{"CREDENTIAL": envResolverMock{in: "test"}}
			iText := tt.in.iText
			iList := tt.in.iList
			iBool := tt.in.iBool
			iPass := tt.in.iPass

			inputManager := NewInputManager(resolvers, iList, iText, iBool, iPass)

			cmd := &exec.Cmd{}
			if tt.in.inType == api.Stdin {
				cmd.Stdin = strings.NewReader(tt.in.stdin)
			}

			got := inputManager.Inputs(cmd, setup, tt.in.inType)

			if got != nil && got.Error() != tt.want.Error() {
				t.Errorf("Inputs(%s) got %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
