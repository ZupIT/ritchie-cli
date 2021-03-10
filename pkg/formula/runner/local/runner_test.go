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

package local

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
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
	ritHome := filepath.Join(tmpDir, ".rit-runner-local")
	repoPath := filepath.Join(ritHome, "repos", "commons")
	defer os.Remove("test.txt")
	defer os.RemoveAll(ritHome)

	inPath := &mocks.InputPathMock{}
	inPath.On("Read", "Type : ").Return("", nil)

	makeBuilder := builder.NewBuildMake()
	batBuilder := builder.NewBuildBat(fileManager)
	shellBuilder := builder.NewBuildShell()

	_ = dirManager.Remove(ritHome)
	_ = dirManager.Remove(repoPath)
	_ = dirManager.Create(repoPath)
	zipFile := filepath.Join("..", "..", "..", "..", "testdata", "ritchie-formulas-test.zip")
	_ = streams.Unzip(zipFile, repoPath)

	envFinder := env.NewFinder(ritHome, fileManager)
	preRunner := NewPreRun(ritHome, makeBuilder, batBuilder, shellBuilder, dirManager, fileManager)

	pInputRunner := prompt.NewInputManager(envResolverMock{in: "test"}, inputMock{}, inputMock{}, inputTextValidatorMock{}, inputTextDefaultMock{}, inputMock{}, inputMock{}, inputMock{}, inPath)
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
		inputResolver formula.InputResolver
		fileManager   stream.FileListMover
		env           env.Finder
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "run local success",
			in: in{
				def:           formula.Definition{Path: "testing/formula", RepoName: "commons"},
				preRun:        preRunner,
				inputResolver: inputResolver,
				fileManager:   fileManager,
				env:           envFinder,
			},
			want: nil,
		},
		{
			name: "success with a non default env",
			in: in{
				def:           formula.Definition{Path: "testing/formula", RepoName: "commons"},
				preRun:        preRunner,
				inputResolver: inputResolver,
				env: envFinderMock{env: env.Holder{
					Current: "prod",
				}},
			},
			want: nil,
		},
		{
			name: "input error local",
			in: in{
				def:           formula.Definition{Path: "testing/formula", RepoName: "commons"},
				preRun:        preRunner,
				inputResolver: inputResolverMock{err: runner.ErrInputNotRecognized},
				fileManager:   fileManager,
				env:           envFinder,
			},
			want: runner.ErrInputNotRecognized,
		},
		{
			name: "pre run error",
			in: in{
				def:           formula.Definition{Path: "testing/formula", RepoName: "commons"},
				preRun:        preRunnerMock{err: errors.New("pre runner error")},
				inputResolver: inputResolver,
				fileManager:   fileManager,
				env:           envFinder,
			},
			want: errors.New("pre runner error"),
		},
		{
			name: "env find error",
			in: in{
				def:           formula.Definition{Path: "testing/formula", RepoName: "commons"},
				preRun:        preRunner,
				inputResolver: inputResolver,
				env:           envFinderMock{err: errors.New("env not found")},
			},
			want: errors.New("env not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in
			local := NewRunner(homeDir, in.fileManager, in.env, in.inputResolver, in.preRun)
			got := local.Run(in.def, api.Prompt, false, nil)

			if tt.want != nil || got != nil {
				assert.EqualError(t, got, tt.want.Error())
			} else {
				assert.FileExists(t, "test.txt")
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

type envResolverMock struct {
	in  string
	err error
}

func (e envResolverMock) Resolve(string) (string, error) {
	return e.in, e.err
}

type inputTextValidatorMock struct{}

func (inputTextValidatorMock) Text(name string, validate func(interface{}) error, helper ...string) (string, error) {
	return "mocked text", nil
}

type inputMock struct {
	text    string
	boolean bool
	items   []string
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

func (i inputMock) Multiselect(formula.Input) ([]string, error) {
	return i.items, i.err
}

type inputTextDefaultMock struct {
	text string
	err  error
}

func (i inputTextDefaultMock) Text(formula.Input) (string, error) {
	return i.text, i.err
}

type envFinderMock struct {
	env env.Holder
	err error
}

func (c envFinderMock) Find() (env.Holder, error) {
	return c.env, c.err
}

type inputResolverMock struct {
	inRunner formula.InputRunner
	err      error
}

func (i inputResolverMock) Resolve(inType api.TermInputType) (formula.InputRunner, error) {
	return i.inRunner, i.err
}
