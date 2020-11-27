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
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

type preRunBuilderTest struct {
	name string

	workspaces        formula.Workspaces
	currentHash       string
	previousHash      string
	currentHashError  error
	previousHashError error
	writeHashError    error

	mustBuild bool
}

var preRunBuilderTests = []preRunBuilderTest{
	{
		name: "should not prompt for rebuild when hash is the same",

		workspaces:        map[string]string{"default": "/pathtodefault"},
		currentHash:       "hash",
		previousHash:      "hash",
		currentHashError:  nil,
		previousHashError: nil,
		writeHashError:    nil,

		mustBuild: false,
	},
	{
		name: "should not prompt for rebuild when hash fails to save",

		workspaces:        map[string]string{"default": "/pathtodefault"},
		currentHash:       "hash",
		previousHash:      "anotherhash",
		currentHashError:  nil,
		previousHashError: nil,
		writeHashError:    fmt.Errorf("Failed to save hash"),

		mustBuild: false,
	},
	{
		name: "should not prompt to rebuild nor fail when no workspaces are returned",

		workspaces:        map[string]string{},
		currentHash:       "hash",
		previousHash:      "anotherhash",
		currentHashError:  nil,
		previousHashError: nil,
		writeHashError:    nil,

		mustBuild: false,
	},
	{
		name: "should not prompt to build when the formula doesn't exist on any workspace",

		workspaces:        map[string]string{"default": "/pathtodefault"},
		currentHash:       "",
		previousHash:      "hash",
		currentHashError:  fmt.Errorf("Formula doesn't exist here"),
		previousHashError: nil,
		writeHashError:    nil,

		mustBuild: false,
	},
	{
		name: "should not prompt to build when no previous hash exists",

		workspaces:        map[string]string{"default": "/pathtodefault"},
		currentHash:       "hash",
		previousHash:      "",
		currentHashError:  nil,
		previousHashError: fmt.Errorf("No previous hash"),
		writeHashError:    nil,

		mustBuild: false,
	},
}

func TestPreRunBuilder(t *testing.T) {
	for _, test := range preRunBuilderTests {
		t.Run(test.name, func(t *testing.T) {
			builderMock := newBuilderMock()

			preRunBuilder := NewPreRunBuilder(workspaceListHasherMock{test.workspaces, test.currentHash, test.currentHashError, test.previousHash,
				test.previousHashError, test.writeHashError}, builderMock)
			preRunBuilder.Build("/testing/formula")

			gotBuilt := builderMock.HasBuilt()
			if gotBuilt != test.mustBuild {
				t.Errorf("Got build %v, wanted %v", gotBuilt, test.mustBuild)
			}
		})
	}
}

type builderMock struct {
	hasBuilt *bool
}

func newBuilderMock() builderMock {
	hasBuilt := false
	return builderMock{&hasBuilt}
}
func (bm builderMock) Build(info formula.BuildInfo) error {
	*bm.hasBuilt = true
	return nil
}
func (bm builderMock) HasBuilt() bool {
	return *bm.hasBuilt
}

type workspaceListHasherMock struct {
	workspaces        formula.Workspaces
	currentHash       string
	currentHashError  error
	previousHash      string
	previousHashError error
	updateHashError   error
}

func (wm workspaceListHasherMock) List() (formula.Workspaces, error) {
	return wm.workspaces, nil
}

func (wm workspaceListHasherMock) CurrentHash(string) (string, error) {
	return wm.currentHash, wm.currentHashError
}

func (wm workspaceListHasherMock) PreviousHash(string) (string, error) {
	return wm.previousHash, wm.previousHashError
}

func (wm workspaceListHasherMock) UpdateHash(string, string) error {
	return wm.updateHashError
}
