/*
 * Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package docker

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/builder"
	"github.com/ZupIT/ritchie-cli/pkg/formula/runner"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/stream/streams"
)

func TestRun(t *testing.T) {
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	tmpDir := os.TempDir()
	homeDir, _ := os.UserHomeDir()
	ritHome := filepath.Join(tmpDir, ".rit-runner-docker")
	repoPath := filepath.Join(ritHome, "repos", "commons")

	dockerBuilder := builder.NewBuildDocker(fileManager)

	_ = dirManager.Remove(ritHome)
	_ = dirManager.Remove(repoPath)
	_ = dirManager.Create(repoPath)
	zipFile := filepath.Join("..", "..", "..", "..", "testdata", "ritchie-formulas-test.zip")
	_ = streams.Unzip(zipFile, repoPath)

	ctxFinder := rcontext.NewFinder(ritHome, fileManager)
	preRunner := NewPreRun(ritHome, dockerBuilder, dirManager, fileManager)
	postRunner := runner.NewPostRunner(fileManager, dirManager)
	inputRunner := runner.NewInput(env.Resolvers{"CREDENTIAL": envResolverMock{in: "test"}}, fileManager, inputMock{}, inputMock{}, inputTextValidatorMock{str: "test"}, inputMock{}, inputMock{})

	type in struct {
		def         formula.Definition
		preRun      formula.PreRunner
		postRun     formula.PostRunner
		inputRun    formula.InputRunner
		context     rcontext.Finder
		fileManager stream.FileWriteExistAppender
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
			name: "run docker success",
			in: in{
				def:         formula.Definition{Path: "testing/formula", RepoName: "commons"},
				preRun:      preRunner,
				postRun:     postRunner,
				inputRun:    inputRunner,
				fileManager: fileManager,
				context:     ctxFinder,
			},
			out: out{
				err: nil,
			},
		},
		{
			name: "input error",
			in: in{
				def:         formula.Definition{Path: "testing/formula", RepoName: "commons"},
				preRun:      preRunner,
				postRun:     postRunner,
				inputRun:    inputRunnerMock{err: runner.ErrInputNotRecognized},
				fileManager: fileManager,
				context:     ctxFinder,
			},
			out: out{
				err: runner.ErrInputNotRecognized,
			},
		},
		{
			name: "pre run error",
			in: in{
				def:         formula.Definition{Path: "testing/formula", RepoName: "commons"},
				preRun:      preRunnerMock{err: errors.New("pre runner error")},
				postRun:     postRunner,
				inputRun:    inputRunner,
				fileManager: fileManager,
				context:     ctxFinder,
			},
			out: out{
				err: errors.New("pre runner error"),
			},
		},
		{
			name: "post run error",
			in: in{
				def:         formula.Definition{Path: "testing/formula", RepoName: "commons"},
				preRun:      preRunner,
				postRun:     postRunnerMock{err: errors.New("post runner error")},
				inputRun:    inputRunner,
				fileManager: fileManager,
				context:     ctxFinder,
			},
			out: out{
				err: nil,
			},
		},
		{
			name: "run docker write .env error",
			in: in{
				def:         formula.Definition{Path: "testing/formula", RepoName: "commons"},
				preRun:      preRunner,
				postRun:     postRunner,
				inputRun:    inputRunner,
				fileManager: fileManagerMock{wErr: errors.New("error to write env file")},
				context:     ctxFinder,
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
				context:     ctxFinder,
			},
			out: out{
				err: errors.New("error to append env file"),
			},
		},
		{
			name: "context find error",
			in: in{
				def:         formula.Definition{Path: "testing/formula", RepoName: "commons"},
				preRun:      preRunner,
				postRun:     postRunner,
				inputRun:    inputRunner,
				fileManager: fileManagerMock{exist: true, aErr: errors.New("error to append env file")},
				context:     ctxFinderMock{err: errors.New("context not found")},
			},
			out: out{
				err: errors.New("context not found"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in
			docker := NewRunner(in.postRun, in.inputRun, in.preRun, in.fileManager, in.context, homeDir)
			got := docker.Run(in.def, api.Prompt, false)

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

func (pr preRunnerMock) PreRun(def formula.Definition) (formula.Setup, error) {
	return pr.setup, pr.err
}

type postRunnerMock struct {
	err error
}

func (po postRunnerMock) PostRun(p formula.Setup, docker bool) error {
	return po.err
}

type envResolverMock struct {
	in  string
	err error
}

func (e envResolverMock) Resolve(string) (string, error) {
	return e.in, e.err
}

type inputTextValidatorMock struct {
	str string
}

func (i inputTextValidatorMock) Text(name string, validate func(interface{}) error, helper ...string) (string, error) {
	return i.str, nil
}

type inputMock struct {
	text    string
	boolean bool
	err     error
}

func (i inputMock) List(string, []string, string, ...string) (string, error) {
	return i.text, i.err
}

func (i inputMock) Text(string, bool, string, ...string) (string, error) {
	return i.text, i.err
}

func (i inputMock) Bool(string, []string, string, ...string) (bool, error) {
	return i.boolean, i.err
}

func (i inputMock) Password(string, ...string) (string, error) {
	return i.text, i.err
}

type ctxFinderMock struct {
	ctx rcontext.ContextHolder
	err error
}

func (c ctxFinderMock) Find() (rcontext.ContextHolder, error) {
	return c.ctx, c.err
}
