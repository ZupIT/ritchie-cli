package runner

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/builder"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/stream/streams"
)

func TestRun(t *testing.T) {
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	tmpDir := os.TempDir()
	ritHome := fmt.Sprintf("%s/.rit-runner", tmpDir)
	repoPath := fmt.Sprintf("%s/repos/commons", ritHome)

	makeBuilder := builder.NewBuildMake()

	_ = dirManager.Remove(ritHome)
	_ = dirManager.Remove(repoPath)
	_ = dirManager.Create(repoPath)
	getwd, _ := os.Getwd()
	fmt.Println(getwd)
	if err := streams.Unzip("../../../testdata/ritchie-formulas-test.zip", repoPath); err != nil {
		t.Error(err)
	}

	preRunner := NewPreRun(ritHome, makeBuilder, dockerBuildMock{}, nil, dirManager, fileManager)
	postRunner := NewPostRunner(fileManager, dirManager)
	inputRunner := NewInput(env.Resolvers{"CREDENTIAL": envResolverMock{in: "test"}}, fileManager, inputMock{}, inputMock{}, inputMock{}, inputMock{})

	type in struct {
		def         formula.Definition
		preRun      formula.PreRunner
		postRun     formula.PostRunner
		inputRun    formula.InputRunner
		fileManager stream.FileWriteExistAppender
		local       bool
	}

	type out struct {
		err error
	}

	tests := []struct {
		name string
		in   in
		out  out
	}{
		{
			name: "Run local success",
			in: in{
				def:         formula.Definition{Path: "testing/formula", RepoName: "commons"},
				preRun:      preRunner,
				postRun:     postRunner,
				inputRun:    inputRunner,
				fileManager: fileManager,
				local:       true,
			},
			out: out{
				err: nil,
			},
		},
		{
			name: "Input error local",
			in: in{
				def:         formula.Definition{Path: "testing/formula", RepoName: "commons"},
				preRun:      preRunner,
				postRun:     postRunner,
				inputRun:    inputRunnerMock{err: ErrInputNotRecognized},
				fileManager: fileManager,
				local:       true,
			},
			out: out{
				err: ErrInputNotRecognized,
			},
		},
		{
			name: "Pre run error",
			in: in{
				def:         formula.Definition{Path: "testing/formula", RepoName: "commons"},
				preRun:      preRunnerMock{err: errors.New("pre runner error")},
				postRun:     postRunner,
				inputRun:    inputRunner,
				fileManager: fileManager,
				local:       true,
			},
			out: out{
				err: errors.New("pre runner error"),
			},
		},
		{
			name: "Post run error",
			in: in{
				def:         formula.Definition{Path: "testing/formula", RepoName: "commons"},
				preRun:      preRunner,
				postRun:     postRunnerMock{err: errors.New("post runner error")},
				inputRun:    inputRunner,
				fileManager: fileManager,
				local:       true,
			},
			out: out{
				err: errors.New("post runner error"),
			},
		},
		{
			name: "Run docker success",
			in: in{
				def:         formula.Definition{Path: "testing/formula", RepoName: "commons"},
				preRun:      preRunner,
				postRun:     postRunner,
				inputRun:    inputRunner,
				fileManager: fileManager,
				local:       false,
			},
			out: out{
				err: nil,
			},
		},
		{
			name: "Input error docker",
			in: in{
				def:         formula.Definition{Path: "testing/formula", RepoName: "commons"},
				preRun:      preRunner,
				postRun:     postRunner,
				inputRun:    inputRunnerMock{err: ErrInputNotRecognized},
				fileManager: fileManager,
				local:       false,
			},
			out: out{
				err: ErrInputNotRecognized,
			},
		},
		{
			name: "Run docker write .env error",
			in: in{
				def:         formula.Definition{Path: "testing/formula", RepoName: "commons"},
				preRun:      preRunner,
				postRun:     postRunner,
				inputRun:    inputRunner,
				fileManager: fileManagerMock{wErr: errors.New("error to write env file")},
				local:       false,
			},
			out: out{
				err: errors.New("error to write env file"),
			},
		},
		{
			name: "Run docker append .env error",
			in: in{
				def:         formula.Definition{Path: "testing/formula", RepoName: "commons"},
				preRun:      preRunner,
				postRun:     postRunner,
				inputRun:    inputRunner,
				fileManager: fileManagerMock{exist: true, aErr: errors.New("error to append env file")},
				local:       false,
			},
			out: out{
				err: errors.New("error to append env file"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in
			runner := NewFormulaRunner(in.postRun, in.inputRun, in.preRun, in.fileManager)
			got := runner.Run(in.def, api.Prompt, in.local)

			if tt.out.err != nil && got != nil && tt.out.err.Error() != got.Error() {
				t.Errorf("Run(%s) got %v, want %v", tt.name, got, tt.out.err)
			}
		})
	}

}

type inputRunnerMock struct {
	err error
}

func (in inputRunnerMock) Inputs(cmd *exec.Cmd, setup formula.Setup, inputType api.TermInputType) error {
	return in.err
}

type preRunnerMock struct {
	setup formula.Setup
	err   error
}

func (pr preRunnerMock) PreRun(def formula.Definition, local bool) (formula.Setup, error) {
	return pr.setup, pr.err
}

type postRunnerMock struct {
	err error
}

func (po postRunnerMock) PostRun(p formula.Setup, docker bool) error {
	return po.err
}
