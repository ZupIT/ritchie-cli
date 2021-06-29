package cmd

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
)

func TestNewAddWorkspaceCmd(t *testing.T) {
	var tests = []struct {
		name          string
		workspaceName string
		workspacePath string
		argsName      string
		argsPath      string
		wantErr       string
		addWSErr      error
		iPathErr      error
	}{
		{
			name:          "add new workspace by prompt",
			workspaceName: "Teste",
			workspacePath: "/home/user/dir",
		},
		{
			name:     "add new workspace by flags",
			argsName: "--name=Teste",
			argsPath: "--path=/home/user/dir",
		},
		{
			name:     "error when one flags is no filled",
			argsName: "--name=Teste",
			argsPath: "",
			wantErr:  "all flags need to be filled",
		},
		{
			name:          "error when workspace does not exists",
			addWSErr:      errors.New("workspace does not exists"),
			workspaceName: "Teste",
			workspacePath: "/home/user/dir",
			wantErr:       "workspace does not exists",
		},
		{
			name:          "error when invalid name",
			workspaceName: "Invalid name",
			workspacePath: "/home/user/dir",
			wantErr:       "the workspace name must not contain spaces",
		},
		{
			name:     "error when invalid name by flags",
			argsName: "--name=Invalid name",
			argsPath: "--path=/home/user/dir",
			wantErr:  "the workspace name must not contain spaces",
		},
		{
			name:          "error when invalid input path",
			workspaceName: "Test",
			iPathErr:      errors.New("input path error"),
			wantErr:       "input path error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			textMock := &mocks.InputTextValidatorMock{}
			textMock.On("Text", "Enter the name of workspace", mock.Anything).Return(tt.workspaceName, tt.wantErr)

			workspaceMock := &mocks.WorkspaceMock{}
			workspaceMock.On("Add", mock.Anything).Return(tt.addWSErr)

			inPath := &mocks.InputPathMock{}
			inPath.On("Read", "Enter the path of workspace (e.g.: /home/user/github) ").Return(tt.workspacePath, tt.iPathErr)

			addNewWorkspace := NewAddWorkspaceCmd(workspaceMock, textMock, inPath)
			addNewWorkspace.SetArgs([]string{tt.argsName, tt.argsPath})

			err := addNewWorkspace.Execute()
			if err != nil {
				assert.Equal(t, tt.wantErr, err.Error())
			} else {
				assert.Empty(t, tt.wantErr)
			}
		})
	}
}
