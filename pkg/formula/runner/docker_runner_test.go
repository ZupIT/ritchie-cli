package runner

import (
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
)

func TestDockerRunner_Run(t *testing.T) {
	def := formula.Definition{
		Path:    "mock/test",
		Bin:     "test-${so}",
		LBin:    "test-${so}",
		MBin:    "test-${so}",
		WBin:    "test-${so}.exe",
		Bundle:  "${so}.zip",
		Config:  "config.json",
		RepoURL: RepoUrl,
	}

	home := os.TempDir()
	_ = fileutil.RemoveDir(home + "/formulas")
	setup := NewDefaultSingleSetup(home, http.DefaultClient)

	type in struct {
		envMock    envResolverMock
		inText     inputMock
		inBool     inputMock
		inPassword inputMock
		preMock    *preRunnerMock
		postMock   *postRunnerMock
		outputR    formula.OutputRunner
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
				inText:  inputMock{text: "ok"},
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
		{
			name: "print and valid error",
			in: in{
				envMock: envResolverMock{in: "ok"},
				inText:  inputMock{text: "ok"},
				inBool:  inputMock{boolean: true},
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
				preRunner = NewDockerPreRunner(setup)
			}

			var postRunner formula.PostRunner
			if in.postMock != nil {
				postRunner = in.postMock
			} else {
				postRunner = NewPostRunner()
			}

			resolvers := env.Resolvers{"test": in.envMock}
			inputManager := NewInputManager(resolvers, in.inText, in.inText, in.inBool, in.inPassword)
			var outputRunner formula.OutputRunner
			if tt.in.outputR == nil {
				outputRunner = NewOutputManager(os.Stdout)
			} else {
				outputRunner = tt.in.outputR
			}
			dockerRunner := NewDockerRunner(preRunner, postRunner, inputManager, outputRunner)

			got := dockerRunner.Run(def, api.Prompt)

			if  (got != nil) != (tt.want != nil) {
				t.Errorf("Run() error = %v, wantErr %v", got, tt.want != nil)
			}

			if got != nil && got.Error() != tt.want.Error() {
				t.Errorf("Run(%s) got %v, want %v", tt.name, got, tt.want)
			}

		})
	}
}

func TestDockerRunner_Run1(t *testing.T) {
	type fields struct {
		PreRunner   formula.PreRunner
		PostRunner  formula.PostRunner
		InputRunner formula.InputRunner
		output      formula.OutputRunner
	}
	type args struct {
		def       formula.Definition
		inputType api.TermInputType
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DockerRunner{
				PreRunner:   tt.fields.PreRunner,
				PostRunner:  tt.fields.PostRunner,
				InputRunner: tt.fields.InputRunner,
				output:      tt.fields.output,
			}
			if err := d.Run(tt.args.def, tt.args.inputType); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}