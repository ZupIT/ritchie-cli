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

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/builder"
	"github.com/ZupIT/ritchie-cli/pkg/formula/input/flag"
	"github.com/ZupIT/ritchie-cli/pkg/formula/input/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/formula/input/stdin"
	"github.com/ZupIT/ritchie-cli/pkg/formula/runner"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/stream/streams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRun(t *testing.T) {
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	tmpDir := os.TempDir()
	homeDir, _ := os.UserHomeDir()
	ritHome := filepath.Join(tmpDir, ".rit-runner-local")
	repoPath := filepath.Join(ritHome, "repos", "commons")
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

	iList := &mocks.InputListMock{}
	iList.On("List", mock.Anything, mock.Anything, mock.Anything).Return("toils", nil)

	iText := &mocks.InputTextMock{}
	iText.On("Text", mock.Anything, mock.Anything, mock.Anything).Return("Test", nil)

	iBool := &mocks.InputBoolMock{}
	iBool.On("Bool", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)

	iPass := &mocks.InputPasswordMock{}
	iPass.On("Password", mock.Anything, mock.Anything, mock.Anything).Return("12345", nil)

	iMultselect := &mocks.InputMultiselectMock{}
	iMultselect.On("Multiselect", mock.Anything).Return([]string{"test", "test"}, nil)

	iTextValidator := &mocks.InputTextValidatorMock{}
	iTextValidator.On("Text")

	iTextDefault := &mocks.InputDefaultTextMock{}
	iTextDefault.On("Text", mock.Anything).Return("test", nil)

	envFinder := env.NewFinder(ritHome, fileManager)

	cf := credential.NewFinder(ritHome, envFinder)
	cs := credential.NewSetter(ritHome, envFinder, dirManager)
	cred := credential.NewResolver(cf, cs, iPass)

	preRunner := NewPreRun(ritHome, makeBuilder, batBuilder, shellBuilder, dirManager, fileManager)
	pInputRunner := prompt.NewInputManager(cred, iList, iText, iTextValidator, iTextDefault,
		iBool, iPass, iMultselect, inPath)
	sInputRunner := stdin.NewInputManager(cred)
	fInputRunner := flag.NewInputManager(cred)

	types := formula.TermInputTypes{
		api.Prompt: pInputRunner,
		api.Stdin:  sInputRunner,
		api.Flag:   fInputRunner,
	}
	inputResolver := runner.NewInputResolver(types)

	tests := []struct {
		name       string
		def        formula.Definition
		inputType  int
		postRunErr error
		envData    env.Holder
		want       error
	}{
		{
			name:      "run local success",
			def:       formula.Definition{Path: "testing/formula", RepoName: "commons"},
			inputType: 0,
			want:      nil,
		},
		{
			name:      "Input error local",
			def:       formula.Definition{Path: "testing/formula", RepoName: "commons"},
			inputType: 3,
			want:      runner.ErrInputNotRecognized,
		},
		{
			name:      "Pre run error",
			def:       formula.Definition{Path: "testing/without-config", RepoName: "commons"},
			inputType: 0,
			want:      errors.New("Failed to load formula config file\nTry running rit update repo\nConfig file path not found: /tmp/.rit-runner-local/repos/commons/testing/without-config/config.json"),
		},
		{
			name:       "Post run error",
			def:        formula.Definition{Path: "testing/formula", RepoName: "commons"},
			postRunErr: errors.New("post runner error"),
			inputType:  0,
			want:       errors.New("post runner error"),
		},
		{
			name:      "env find error",
			def:       formula.Definition{Path: "testing/formula", RepoName: "commons"},
			inputType: 0,
			envData:   env.Holder{},
			want:      errors.New("env not found"),
		},
		{
			name:      "success with a non default env",
			def:       formula.Definition{Path: "testing/formula", RepoName: "commons"},
			inputType: 0,
			envData:   env.Holder{Current: "prd"},
			want:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eMock := &mocks.EnvFinderMock{}
			eMock.On("Find").Return(tt.envData, tt.want)

			pr := &mocks.PostRunner{}
			pr.On("PostRun", mock.Anything, mock.Anything).Return(tt.postRunErr)

			local := NewRunner(pr, inputResolver, preRunner, fileManager, eMock, homeDir)
			got := local.Run(tt.def, api.TermInputType(tt.inputType), false, nil)

			assert.Equal(t, tt.want, got)
		})
	}

}
