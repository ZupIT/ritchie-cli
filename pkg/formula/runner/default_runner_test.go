package runner

import (
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

var RepoUrl = os.Getenv("REPO_URL")

func TestDefaultRunner_Run(t *testing.T) {
	def := formula.Definition{
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
		inPass   inputMock
		preMock  *preRunnerMock
		postMock *postRunnerMock
		outputR  formula.OutputRunner
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
				inPass:  inputMock{text: "******"},
			},
			want: nil,
		},
		{
			name: "pre run error",
			in: in{
				envMock: envResolverMock{in: "ok"},
				inText:  inputMock{text: "ok"},
				inBool:  inputMock{boolean: true},
				inPass:  inputMock{text: "******"},
				preMock: &preRunnerMock{
					setup: formula.Setup{},
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
				inPass:  inputMock{text: "******"},
			},
			want: errors.New("fail to resolve input"),
		},
		{
			name: "post run error",
			in: in{
				envMock:  envResolverMock{in: "ok"},
				inText:   inputMock{text: "ok"},
				inBool:   inputMock{boolean: true},
				inPass:   inputMock{text: "******"},
				postMock: &postRunnerMock{error: errors.New("error in remove dir")},
			},
			want: errors.New("error in remove dir"),
		},
		{
			name: "print and valid error",
			in: in{
				envMock: envResolverMock{in: "ok"},
				inText:  inputMock{text: ""},
				inBool:  inputMock{boolean: true},
				inPass:  inputMock{text: "******"},
				outputR: outputMock{
					validAndPrint: func(setup formula.Setup) error {
						return errors.New("some Error on output")
					},
				},
			},
			want: errors.New("some Error on output"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in

			var preRunner formula.PreRunner
			if in.preMock != nil {
				preRunner = in.preMock
			} else {
				preRunner = NewDefaultPreRunner(setup)
			}

			var postRunner formula.PostRunner
			if in.postMock != nil {
				postRunner = in.postMock
			} else {
				postRunner = NewPostRunner()
			}

			resolvers := env.Resolvers{"test": in.envMock}
			inputManager := NewInputManager(resolvers, in.inText, in.inText, in.inBool, in.inPass)
			var outputRunner formula.OutputRunner
			if tt.in.outputR == nil {
				outputRunner = NewOutputManager(os.Stdout)
			} else {
				outputRunner = tt.in.outputR
			}
			defaultRunner := NewDefaultRunner(preRunner, postRunner, inputManager, outputRunner)

			got := defaultRunner.Run(def, api.Prompt)

			if (got != nil) != (tt.want != nil) {
				t.Errorf("Run() error = %v, wantErr %v", got, tt.want != nil)
			}

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

func (i inputMock) TextWithValidate(name string, validate func(string) error) (string, error) {
	return i.text, i.err
}

func (i inputMock) Bool(string, []string) (bool, error) {
	return i.boolean, i.err
}

func (i inputMock) Password(string) (string, error) {
	return i.text, i.err
}

type outputMock struct {
	validAndPrint func(setup formula.Setup) error
}

func (o outputMock) ValidAndPrint(setup formula.Setup) error {
	return o.validAndPrint(setup)
}

type envResolverMock struct {
	in  string
	err error
}

func (e envResolverMock) Resolve(string) (string, error) {
	return e.in, e.err
}

type preRunnerMock struct {
	setup formula.Setup
	error error
}

func (p preRunnerMock) PreRun(formula.Definition) (formula.Setup, error) {
	return p.setup, p.error
}

type postRunnerMock struct {
	error error
}

func (p postRunnerMock) PostRun(formula.Setup, bool) error {
	return p.error
}
