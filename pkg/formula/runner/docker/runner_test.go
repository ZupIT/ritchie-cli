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
	"path/filepath"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/builder"
	"github.com/ZupIT/ritchie-cli/pkg/formula/input/flag"
	"github.com/ZupIT/ritchie-cli/pkg/formula/input/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/formula/input/stdin"
	"github.com/ZupIT/ritchie-cli/pkg/formula/runner"
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

	ctxFinder := env.NewFinder(ritHome, fileManager)
	preRunner := NewPreRun(ritHome, dockerBuilder, dirManager, fileManager)
	postRunner := runner.NewPostRunner(fileManager, dirManager)
	pInputRunner := prompt.NewInputManager(envResolverMock{in: "test"}, fileManager, inputMock{}, inputMock{}, inputTextValidatorMock{str: "test"}, inputTextDefaultMock{}, inputMock{}, inputMock{})
	sInputRunner := stdin.NewInputManager(envResolverMock{in: "test"})
	fInputRunner := flag.NewInputManager(envResolverMock{in: "test"})

	types := formula.TermInputTypes{
		api.Prompt: pInputRunner,
		api.Stdin:  sInputRunner,
		api.Flag:   fInputRunner,
	}
	inputResolver := runner.NewInputResolver(types)

	type in struct {
		def           formula.Definition
		preRun        formula.PreRunner
		postRun       formula.PostRunner
		inputResolver formula.InputResolver
		context       env.Finder
		fileManager   stream.FileWriteExistAppender
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
				def:           formula.Definition{Path: "testing/formula", RepoName: "commons"},
				preRun:        preRunner,
				postRun:       postRunner,
				inputResolver: inputResolver,
				fileManager:   fileManager,
				context:       ctxFinder,
			},
			out: out{
				err: nil,
			},
		},
		{
			name: "input error",
			in: in{
				def:           formula.Definition{Path: "testing/formula", RepoName: "commons"},
				preRun:        preRunner,
				postRun:       postRunner,
				inputResolver: inputResolverMock{err: runner.ErrInputNotRecognized},
				fileManager:   fileManager,
				context:       ctxFinder,
			},
			out: out{
				err: runner.ErrInputNotRecognized,
			},
		},
		{
			name: "pre run error",
			in: in{
				def:           formula.Definition{Path: "testing/formula", RepoName: "commons"},
				preRun:        preRunnerMock{err: errors.New("pre runner error")},
				postRun:       postRunner,
				inputResolver: inputResolver,
				fileManager:   fileManager,
				context:       ctxFinder,
			},
			out: out{
				err: errors.New("pre runner error"),
			},
		},
		{
			name: "post run error",
			in: in{
				def:           formula.Definition{Path: "testing/formula", RepoName: "commons"},
				preRun:        preRunner,
				postRun:       postRunnerMock{err: errors.New("post runner error")},
				inputResolver: inputResolver,
				fileManager:   fileManager,
				context:       ctxFinder,
			},
			out: out{
				err: nil,
			},
		},
		{
			name: "run docker write .env error",
			in: in{
				def:           formula.Definition{Path: "testing/formula", RepoName: "commons"},
				preRun:        preRunner,
				postRun:       postRunner,
				inputResolver: inputResolver,
				fileManager:   fileManagerMock{wErr: errors.New("error to write env file")},
				context:       ctxFinder,
			},
			out: out{
				err: errors.New("error to write env file"),
			},
		},
		{
			name: "Run docker append .env error",
			in: in{
				def:           formula.Definition{Path: "testing/formula", RepoName: "commons"},
				preRun:        preRunner,
				postRun:       postRunner,
				inputResolver: inputResolver,
				fileManager:   fileManagerMock{exist: true, aErr: errors.New("error to append env file")},
				context:       ctxFinder,
			},
			out: out{
				err: errors.New("error to append env file"),
			},
		},
		{
			name: "context find error",
			in: in{
				def:           formula.Definition{Path: "testing/formula", RepoName: "commons"},
				preRun:        preRunner,
				postRun:       postRunner,
				inputResolver: inputResolver,
				fileManager:   fileManagerMock{exist: true, aErr: errors.New("error to append env file")},
				context:       ctxFinderMock{err: errors.New("context not found")},
			},
			out: out{
				err: errors.New("context not found"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in
			docker := NewRunner(in.postRun, in.inputResolver, in.preRun, in.fileManager, in.context, homeDir)
			got := docker.Run(in.def, api.Prompt, false, nil)

			if tt.out.err != nil && got != nil && tt.out.err.Error() != got.Error() {
				t.Errorf("Run(%s) got %v, want %v", tt.name, got, tt.out.err)
			}
		})
	}

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

func (i inputMock) List(string, []string, ...string) (string, error) {
	return i.text, i.err
}

func (i inputMock) Text(string, bool, ...string) (string, error) {
	return i.text, i.err
}

func (i inputMock) Bool(string, []string, ...string) (bool, error) {
	return i.boolean, i.err
}

func (i inputMock) Password(string, ...string) (string, error) {
	return i.text, i.err
}

type inputTextDefaultMock struct {
	text string
	err  error
}

func (i inputTextDefaultMock) Text(formula.Input) (string, error) {
	return i.text, i.err
}

type ctxFinderMock struct {
	ctx env.Holder
	err error
}

func (c ctxFinderMock) Find() (env.Holder, error) {
	return c.ctx, c.err
}

type inputResolverMock struct {
	inRunner formula.InputRunner
	err      error
}

func (i inputResolverMock) Resolve(inType api.TermInputType) (formula.InputRunner, error) {
	return i.inRunner, i.err
}
