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

package runner

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
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/stream/streams"
)

func TestRun(t *testing.T) {
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	tmpDir := os.TempDir()
	homeDir, _ := os.UserHomeDir()
	ritHome := filepath.Join(tmpDir, ".rit-runner")
	repoPath := filepath.Join(ritHome, "repos", "commons")

	makeBuilder := builder.NewBuildMake()
	batBuilder := builder.NewBuildBat(fileManager)

	_ = dirManager.Remove(ritHome)
	_ = dirManager.Remove(repoPath)
	_ = dirManager.Create(repoPath)
	zipFile := filepath.Join("..", "..", "..", "testdata", "ritchie-formulas-test.zip")
	_ = streams.Unzip(zipFile, repoPath)

	ctxFinder := rcontext.NewFinder(ritHome, fileManager)
	preRunner := NewPreRun(ritHome, makeBuilder, dockerBuildMock{}, batBuilder, dirManager, fileManager)
	postRunner := NewPostRunner(fileManager, dirManager)
	inputRunner := NewInput(env.Resolvers{"CREDENTIAL": envResolverMock{in: "test"}}, fileManager, inputMock{}, inputMock{}, inputMock{}, inputMock{})

	type in struct {
		def         formula.Definition
		preRun      formula.PreRunner
		postRun     formula.PostRunner
		inputRun    formula.InputRunner
		fileManager stream.FileWriteExistAppender
		docker      bool
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
			name: "run local success",
			in: in{
				def:         formula.Definition{Path: "testing/formula", RepoName: "commons"},
				preRun:      preRunner,
				postRun:     postRunner,
				inputRun:    inputRunner,
				fileManager: fileManager,
				docker:      false,
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
				docker:      false,
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
				docker:      false,
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
				docker:      false,
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
				docker:      true,
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
				docker:      true,
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
				docker:      true,
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
				docker:      true,
			},
			out: out{
				err: errors.New("error to append env file"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in
			runner := NewFormulaRunner(in.postRun, in.inputRun, in.preRun, in.fileManager, ctxFinder)
			got := runner.Run(in.def, api.Prompt, in.docker, false, homeDir)

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
