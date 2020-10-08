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
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

type preRunBuilderTest struct {
	name string

	workspaces          formula.Workspaces
	rebuildPromptAnswer bool
	currentHash         string
	previousHash        string
	createHashDirError  error
	writeHashError      error
	promptError         error

	mustBuild         bool
	mustPromptRebuild bool
}

var preRunBuilderTests = []preRunBuilderTest{
	{
		name: "should not prompt for rebuild when hash is the same",

		workspaces:          map[string]string{"default": "/pathtodefault"},
		rebuildPromptAnswer: false,
		currentHash:         "hash",
		previousHash:        "hash",
		createHashDirError:  nil,
		writeHashError:      nil,
		promptError:         nil,

		mustBuild:         false,
		mustPromptRebuild: false,
	},
	{
		name: "should rebuild when user chooses to",

		workspaces:          map[string]string{"default": "/pathtodefault"},
		rebuildPromptAnswer: true,
		currentHash:         "hash",
		previousHash:        "anotherhash",
		createHashDirError:  nil,
		writeHashError:      nil,
		promptError:         nil,

		mustBuild:         true,
		mustPromptRebuild: true,
	},
	{
		name: "should not rebuild when user chooses not to",

		workspaces:          map[string]string{"default": "/pathtodefault"},
		rebuildPromptAnswer: false,
		currentHash:         "hash",
		previousHash:        "anotherhash",
		createHashDirError:  nil,
		writeHashError:      nil,
		promptError:         nil,

		mustBuild:         false,
		mustPromptRebuild: true,
	},
	{
		name: "should not prompt for rebuild when hash fails to save",

		workspaces:          map[string]string{"default": "/pathtodefault"},
		rebuildPromptAnswer: false,
		currentHash:         "hash",
		previousHash:        "anotherhash",
		createHashDirError:  nil,
		writeHashError:      fmt.Errorf("Failed to save hash"),
		promptError:         nil,

		mustBuild:         false,
		mustPromptRebuild: false,
	},
	{
		name: "should ignore directory creation errors",

		workspaces:          map[string]string{"default": "/pathtodefault"},
		rebuildPromptAnswer: false,
		currentHash:         "hash",
		previousHash:        "anotherhash",
		createHashDirError:  fmt.Errorf("Failed to create dir"),
		writeHashError:      nil,
		promptError:         nil,

		mustBuild:         false,
		mustPromptRebuild: true,
	},
	{
		name: "should not prompt to rebuild nor fail when no workspaces are returned",

		workspaces:          map[string]string{},
		rebuildPromptAnswer: true,
		currentHash:         "hash",
		previousHash:        "anotherhash",
		createHashDirError:  nil,
		writeHashError:      nil,
		promptError:         nil,

		mustBuild:         false,
		mustPromptRebuild: false,
	},
	{
		name: "should not build when user Ctrl+C's on prompt",

		workspaces:          map[string]string{"default": "/pathtodefault"},
		rebuildPromptAnswer: true,
		currentHash:         "hash",
		previousHash:        "anotherhash",
		createHashDirError:  nil,
		writeHashError:      nil,
		promptError:         fmt.Errorf("Ctrl+C on survey"),

		mustBuild:         false,
		mustPromptRebuild: true,
	},
}

func TestPreRunBuilder(t *testing.T) {
	tmpDir := os.TempDir()
	ritHomeName := ".rit-pre-run-builder"
	ritHome := filepath.Join(tmpDir, ritHomeName)

	for _, test := range preRunBuilderTests {
		t.Run(test.name, func(t *testing.T) {
			builderMock := newBuilderMock()
			inputBoolMock := newInputBoolMock(test.rebuildPromptAnswer, test.promptError)

			preRunBuilder := NewPreRunBuilder(ritHome, workspaceListerMock{test.workspaces}, builderMock,
				dirHashManagerMock{test.createHashDirError, nil, test.currentHash, nil},
				fileManagerMock{[]byte(test.previousHash), nil, test.writeHashError, nil, nil, nil, nil, false, []string{}},
				inputBoolMock)
			preRunBuilder.Build("/testing/formula")

			if builderMock.HasBuilt() != test.mustBuild {
				if test.mustBuild {
					t.Error("Expected formula to build but it didn't")
				} else {
					t.Error("Expected formula not to build but it did")
				}
			}

			if inputBoolMock.HasBeenCalled() != test.mustPromptRebuild {
				if test.mustPromptRebuild {
					t.Error("Expected formula to prompt for rebuild but it didn't")
				} else {
					t.Error("Expected formula not to prompt for rebuild but it did")
				}
			}
		})
	}
}

type inputBoolMock struct {
	hasBeenCalled *bool
	answer        bool
	err           error
}

func newInputBoolMock(answer bool, err error) inputBoolMock {
	hasBeenCalled := false
	return inputBoolMock{&hasBeenCalled, answer, err}
}
func (in inputBoolMock) Bool(string, []string, ...string) (bool, error) {
	*in.hasBeenCalled = true
	return in.answer, in.err
}
func (in inputBoolMock) HasBeenCalled() bool {
	return *in.hasBeenCalled
}

type builderMock struct {
	hasBuilt *bool
}

func newBuilderMock() builderMock {
	hasBuilt := false
	return builderMock{&hasBuilt}
}
func (bm builderMock) Build(string, string) error {
	*bm.hasBuilt = true
	return nil
}
func (bm builderMock) HasBuilt() bool {
	return *bm.hasBuilt
}

type dirHashManagerMock struct {
	createErr error
	removeErr error
	hash      string
	hashErr   error
}

func (di dirHashManagerMock) Create(dir string) error {
	return di.createErr
}
func (di dirHashManagerMock) Remove(dir string) error {
	return di.removeErr
}
func (di dirHashManagerMock) Hash(dir string) (string, error) {
	return di.hash, di.hashErr
}

type workspaceListerMock struct {
	workspaces formula.Workspaces
}

func (wm workspaceListerMock) List() (formula.Workspaces, error) {
	return wm.workspaces, nil
}
