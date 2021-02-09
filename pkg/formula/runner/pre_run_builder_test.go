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
	"fmt"
	"testing"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type preRunBuilderTest struct {
	name string

	workspaces        formula.Workspaces
	currentHash       string
	previousHash      string
	currentHashError  error
	previousHashError error
	writeHashError    error
	builderError      error
	listError         error

	errExpected error
}

var preRunBuilderTests = []preRunBuilderTest{
	{
		name: "should not prompt for rebuild when hash is the same",

		workspaces:   map[string]string{"default": "/pathtodefault"},
		currentHash:  "hash",
		previousHash: "hash",
	},
	{
		name: "return error when hash fails to save",

		workspaces:     map[string]string{"default": "/pathtodefault"},
		currentHash:    "hash",
		previousHash:   "anotherhash",
		writeHashError: fmt.Errorf("Failed to save hash"),

		errExpected: errors.New("Failed to detect formula changes, executing the last build: Failed to save hash"),
	},
	{
		name: "return nil when no workspaces are returned",

		workspaces:   map[string]string{},
		currentHash:  "hash",
		previousHash: "anotherhash",
	},
	{
		name: "return nil when the formula doesn't exist on any workspace",

		workspaces:       map[string]string{"default": "/pathtodefault"},
		currentHash:      "",
		previousHash:     "hash",
		currentHashError: fmt.Errorf("Formula doesn't exist here"),
	},
	{
		name: "return nil when no previous hash exists",

		workspaces:        map[string]string{"default": "/pathtodefault"},
		currentHash:       "hash",
		previousHash:      "",
		previousHashError: fmt.Errorf("No previous hash"),
	},
	{
		name: "return error when build workspace returns error",

		workspaces:   map[string]string{"default": "/pathtodefault"},
		currentHash:  "hashtwo",
		previousHash: "hash",
		builderError: fmt.Errorf("Some error builder"),

		errExpected: errors.New("Failed to build formula: Some error builder"),
	},
	{
		name: "returns null when routine runs without errors",

		workspaces:   map[string]string{"default": "/pathtodefault"},
		currentHash:  "hashtwo",
		previousHash: "hash",
	},
	{
		name: "returns error when list workspace returns error",

		workspaces:   map[string]string{"default": "/pathtodefault"},
		currentHash:  "hashtwo",
		previousHash: "hash",
		listError:    fmt.Errorf("Some error list"),

		errExpected: errors.New("Failed to detect formula changes, executing the last build: Some error list"),
	},
}

func TestPreRunBuilder(t *testing.T) {
	for _, test := range preRunBuilderTests {
		t.Run(test.name, func(t *testing.T) {
			builderMock := new(mocks.BuilderMock)
			builderMock.On("Build", mock.Anything).Return(test.builderError)
			builderMock.On("HasBuilt", mock.Anything).Return(false)

			workspaceListHasherMock := new(mocks.WorkspaceListHasherMock)
			workspaceListHasherMock.On("List").Return(test.workspaces, test.listError)
			workspaceListHasherMock.On("CurrentHash", mock.Anything).Return(test.currentHash, test.currentHashError)
			workspaceListHasherMock.On("PreviousHash", mock.Anything).Return(test.previousHash, test.previousHashError)
			workspaceListHasherMock.On("UpdateHash", mock.Anything, mock.Anything).Return(test.writeHashError)

			preRunBuilder := NewPreRunBuilder(workspaceListHasherMock, builderMock)

			got := preRunBuilder.Build("/testing/formula")
			assert.Equal(t, test.errExpected, got)
		})
	}
}
