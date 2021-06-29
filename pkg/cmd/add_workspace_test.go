package cmd

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/formula/workspace"
)

func TestNewAddWorkspaceCmd(t *testing.T) {
	var tests = []struct {
		name          string
		workspaceName string
		workspacePath string
		argsName      string
		argsPath      string
		want          error
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
			want:     errors.New("all flags need to be filled"),
		},
		{
			name:          "error when workspace does not exists",
			addWSErr:      errors.New("workspace does not exists"),
			workspaceName: "Teste",
			workspacePath: "/home/user/dir",
			want:          errors.New("workspace does not exists"),
		},
		{
			name:          "error when invalid name",
			workspaceName: "Invalid name",
			workspacePath: "/home/user/dir",
			want:          workspace.ErrInvalidWorkspaceName,
		},
		{
			name:     "error when invalid name by flags",
			argsName: "--name=Invalid name",
			argsPath: "--path=/home/user/dir",
			want:     workspace.ErrInvalidWorkspaceName,
		},
		{
			name:          "error when invalid input path",
			workspaceName: "Test",
			iPathErr:      errors.New("input path error"),
			want:          errors.New("input path error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			textMock := &mocks.InputTextValidatorMock{}
			textMock.On("Text", "Enter the name of workspace", mock.Anything).Return(tt.workspaceName, tt.want)

			workspaceMock := &mocks.WorkspaceMock{}
			workspaceMock.On("Add", mock.Anything).Return(tt.addWSErr)

			inPath := &mocks.InputPathMock{}
			inPath.On("Read", "Enter the path of workspace (e.g.: /home/user/github) ").Return(tt.workspacePath, tt.iPathErr)

			addNewWorkspace := NewAddWorkspaceCmd(workspaceMock, textMock, inPath)
			addNewWorkspace.SetArgs([]string{tt.argsName, tt.argsPath})

			got := addNewWorkspace.Execute()
			if got != nil {
				assert.EqualError(t, got, tt.want.Error())
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
