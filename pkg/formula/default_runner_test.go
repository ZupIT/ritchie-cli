package formula

import (
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
)

func TestDefaultRunner_Run(t *testing.T) {
	def := Definition{
		Path:    "mock/test",
		Bin:     "test-linux",
		Bundle:  "linux.zip",
		Config:  "config.json",
		RepoURL: RepoUrl,
	}

	home := os.TempDir()
	_ = fileutil.RemoveDir(home + "/formulas")
	setup := NewDefaultSingleSetup(home, http.DefaultClient)

	type in struct {
		envMock  envResolverMock
		inText   inputMock
		inBool   inputMock
		preMock  *preRunnerMock
		postMock *postRunnerMock
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "success",
			in: in{
				envMock: envResolverMock{in: "ok"},
				inText:  inputMock{text: ""},
				inBool:  inputMock{boolean: true},
			},
			want: nil,
		},
		{
			name: "pre run error",
			in: in{
				envMock: envResolverMock{in: "ok"},
				inText:  inputMock{text: "ok"},
				inBool:  inputMock{boolean: true},
				preMock: &preRunnerMock{
					setup: Setup{},
					error: ErrFormulaBinNotFound,
				},
			},
			want: ErrFormulaBinNotFound,
		},
		{
			name: "inputs error",
			in: in{
				envMock: envResolverMock{in: "ok"},
				inText:  inputMock{err: errors.New("fail to resolve input")},
				inBool:  inputMock{boolean: true},
			},
			want: errors.New("fail to resolve input"),
		},
		{
			name: "post run error",
			in: in{
				envMock:  envResolverMock{in: "ok"},
				inText:   inputMock{text: "ok"},
				inBool:   inputMock{boolean: true},
				postMock: &postRunnerMock{error: errors.New("error in remove dir")},
			},
			want: errors.New("error in remove dir"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in

			var preRunner PreRunner
			if in.preMock != nil {
				preRunner = in.preMock
			} else {
				preRunner = NewDefaultPreRunner(setup)
			}

			var postRunner PostRunner
			if in.postMock != nil {
				postRunner = in.postMock
			} else {
				postRunner = NewPostRunner()
			}

			resolvers := env.Resolvers{"test": in.envMock}
			inputManager := NewInputManager(resolvers, in.inText, in.inText, in.inBool)
			defaultRunner := NewDefaultRunner(preRunner, postRunner, inputManager)

			got := defaultRunner.Run(def, api.Prompt)

			if got != nil && got.Error() != tt.want.Error() {
				t.Errorf("Run(%s) got %v, want %v", tt.name, got, tt.want)
			}

		})
	}
}

type inputMock struct {
	text    string
	boolean bool
	err     error
}

func (i inputMock) List(string, []string) (string, error) {
	return i.text, i.err
}

func (i inputMock) Text(string, bool) (string, error) {
	return i.text, i.err
}

func (i inputMock) Bool(string, []string) (bool, error) {
	return i.boolean, i.err
}

type envResolverMock struct {
	in  string
	err error
}

func (e envResolverMock) Resolve(string) (string, error) {
	return e.in, e.err
}

type preRunnerMock struct {
	setup Setup
	error error
}

func (p preRunnerMock) PreRun(Definition) (Setup, error) {
	return p.setup, p.error
}

type postRunnerMock struct {
	error error
}

func (p postRunnerMock) PostRun(Setup, bool) error {
	return p.error
}
