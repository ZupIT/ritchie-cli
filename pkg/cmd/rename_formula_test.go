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

package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

func TestRenameFormulaCmd(t *testing.T) {
	type in struct {
		inputText     string
		inputTextErr  error
		inputTextVal  string
		inputList     string
		inputListErr  error
		wspaceList    formula.Workspaces
		wspaceListErr error
		wspaceAddErr  error
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "success",
			in: in{
				inputTextVal: "rit test test",
				inputList:    "go",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspaceMock := new(mocks.WorkspaceForm)
			workspaceMock.On("List").Return(tt.in.wspaceList, tt.in.wspaceListErr)
			workspaceMock.On("Add", mock.Anything).Return(tt.in.wspaceAddErr)
			workspaceMock.On("CurrentHash", mock.Anything).Return("48d47029-2abf-4a2e-b5f2-f5b60471423e", nil)
			workspaceMock.On("UpdateHash", mock.Anything, mock.Anything).Return(nil)

			inputTextMock := new(mocks.InputTextMock)
			inputTextMock.On("Text", mock.Anything, mock.Anything, mock.Anything).Return(tt.in.inputText, tt.in.inputTextErr)

			inputListMock := new(mocks.InputListMock)
			inputListMock.On("List", mock.Anything, mock.Anything, mock.Anything).Return(tt.in.inputList, tt.in.inputListErr)

			inPath := &mocks.InputPathMock{}
			inPath.On("Read", "Workspace path (e.g.: /home/user/github): ").Return("", nil)

			renameFormulaCmd := NewRenameFormulaCmd(
				workspaceMock,
				inputTextMock,
				inputListMock,
				inPath,
			)

			got := renameFormulaCmd.Execute()

			assert.Equal(t, tt.want, got)
		})
	}

}
