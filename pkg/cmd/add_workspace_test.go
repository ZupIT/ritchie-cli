package cmd

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

func TestNewAddWorkspaceCmd(t *testing.T) {
	var tests = []struct {
		name                  string
		addWorkspaceWithError error
		workspaceName         string
		workspacePath         string
		argsName              string
		argsPath              string
		wantErr               string
	}{
		{
			name:                  "add new workspace by prompt",
			addWorkspaceWithError: nil,
			workspaceName:         "Teste",
			workspacePath:         "/home/user/dir",
			wantErr:               "",
		},
		{
			name:                  "add new workspace by flags",
			addWorkspaceWithError: nil,
			workspaceName:         "Teste",
			workspacePath:         "/home/user/dir",
			argsName:              "--name=Teste",
			argsPath:              "--path=/home/user/dir",
			wantErr:               "",
		},
		{
			name:                  "error when one flags is no filled",
			addWorkspaceWithError: nil,
			workspaceName:         "Teste",
			workspacePath:         "",
			argsName:              "--name=Teste",
			argsPath:              "",
			wantErr:               "all flags need to be filled",
		},
		{
			name:                  "error when workspace does not exists",
			addWorkspaceWithError: errors.New("workspace does not exists"),
			workspaceName:         "Teste",
			workspacePath:         "/home/user/dir",
			wantErr:               "workspace does not exists",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			textMock := &mocks.InputTextMock{}
			textMock.On("Text", "Enter the name of workspace", true, mock.Anything).Return(tt.workspaceName, nil)
			textMock.On("Text", "Enter the path of workspace (e.g.: /home/user/github)", true, mock.Anything).Return(tt.workspacePath, nil)

			wspace := formula.Workspace{
				Name: tt.workspaceName,
				Dir:  tt.workspacePath,
			}

			workspaceMock := &mocks.WorkspaceMock{}
			workspaceMock.On("Add", wspace).Return(tt.addWorkspaceWithError)

			inPath := mocks.InputPathMock{}

			addNewWorkspace := NewAddWorkspaceCmd(workspaceMock, textMock, inPath)
			addNewWorkspace.SetArgs([]string{tt.argsName, tt.argsPath})

			err := addNewWorkspace.Execute()
			if err != nil {
				assert.Equal(t, err.Error(), tt.wantErr)
			} else {
				assert.Empty(t, tt.wantErr)
			}
		})
	}
}
