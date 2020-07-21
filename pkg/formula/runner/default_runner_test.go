package runner

import (
	"errors"
	"runtime"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/os/osutil"
)

func TestDefaultRunner_Run(t *testing.T) {
	def := formula.Definition{
		Path: "mock/test",
	}

	type in struct {
		envMock  envResolverMock
		inText   inputMock
		inBool   inputMock
		inPass   inputMock
		preMock  *preRunnerMock
		postMock *postRunnerMock
	}

	var binName string
	switch runtime.GOOS {
	case osutil.Windows:
		binName = "../../../testdata/run-mock.bat"
	default:
		binName = "../../../testdata/run-mock.sh"
	}
	defaultPreMock := &preRunnerMock{
		setup: formula.Setup{
			BinName: binName,
			Config: formula.Config{
				Inputs: []formula.Input{
					{
						Name: "SOME_INPUT",
						Type: "text",
					},
				},
			},
		},
	}

	defaultPostMock := &postRunnerMock{}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "success",
			in: in{
				envMock:  envResolverMock{in: "ok"},
				inText:   inputMock{text: ""},
				inBool:   inputMock{boolean: true},
				inPass:   inputMock{text: "******"},
				preMock:  defaultPreMock,
				postMock: defaultPostMock,
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
				envMock:  envResolverMock{in: "ok"},
				inText:   inputMock{err: errors.New("fail to resolve input")},
				inBool:   inputMock{boolean: true},
				inPass:   inputMock{text: "******"},
				preMock:  defaultPreMock,
				postMock: defaultPostMock,
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
				preMock:  defaultPreMock,
				postMock: &postRunnerMock{error: errors.New("error in remove dir")},
			},
			want: errors.New("error in remove dir"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in

			resolvers := env.Resolvers{"test": in.envMock}
			inputManager := NewInputManager(resolvers, in.inText, in.inText, in.inBool, in.inPass)
			defaultRunner := NewDefaultRunner(in.preMock, in.postMock, inputManager)

			got := defaultRunner.Run(def, api.Prompt)

			if (tt.want != nil && got == nil) || got != nil && got.Error() != tt.want.Error() {
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

func (i inputMock) Text(string, bool, ...string) (string, error) {
	return i.text, i.err
}

func (i inputMock) TextWithValidate(string, func(interface{}) error, ...string) (string, error) {
	return i.text, i.err
}

func (i inputMock) Bool(string, []string) (bool, error) {
	return i.boolean, i.err
}

func (i inputMock) Password(string) (string, error) {
	return i.text, i.err
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
