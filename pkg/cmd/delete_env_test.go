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
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewDeleteEnv(t *testing.T) {
	tmpDir := os.TempDir()
	ritHomeDir := filepath.Join(tmpDir, ".rit")
	envFile := filepath.Join(ritHomeDir, env.FileName)
	_ = os.MkdirAll(ritHomeDir, os.ModePerm)
	defer os.RemoveAll(ritHomeDir)

	fileManager := stream.NewFileManager()

	envFinder := env.NewFinder(ritHomeDir, fileManager)
	envRemover := env.NewRemover(ritHomeDir, envFinder, fileManager)
	envFindRemover := env.NewFindRemover(envFinder, envRemover)

	envEmpty := env.Holder{Current: "", All: []string{}}
	envCompleted := env.Holder{Current: "prod", All: []string{"prod", "qa", "stg"}}

	tests := []struct {
		name            string
		env             env.Holder
		envFileNil      bool
		inputBoolError  error
		inputBoolResult bool
		inputListString string
		inputListError  error
		wantErr         string
		envResultInFile env.Holder
		args            []string
	}{
		{
			name:            "execute prompt with success",
			inputBoolResult: true,
			inputListString: "qa",
			env:             envCompleted,
			envResultInFile: env.Holder{Current: "prod", All: []string{"prod", "stg"}},
		},
		{
			name:            "execute flag with success",
			inputListString: "prod",
			env:             envCompleted,
			envResultInFile: env.Holder{Current: "", All: []string{"qa", "stg"}},
			args:            []string{"--env=prod"},
		},
		{
			name:            "execute with success when not envs defined",
			env:             envEmpty,
			envResultInFile: envEmpty,
		},
		{
			name:           "fail on input list error",
			inputListError: errors.New("some error"),
			env:            envCompleted,
			wantErr:        "some error",
		},
		{
			name:            "fail on input bool error",
			inputBoolError:  errors.New("some error"),
			env:             envCompleted,
			envResultInFile: envEmpty,
			wantErr:         "some error",
		},
		{
			name:            "fail on flag error when flag not defined",
			inputListString: "prod",
			env:             envCompleted,
			envResultInFile: env.Holder{Current: "", All: []string{"qa", "stg"}},
			args:            []string{"--list=nana"},
			wantErr:         "unknown flag: --list",
		},
		{
			name:            "fail on flag error when env defined nil",
			env:             envCompleted,
			envResultInFile: env.Holder{Current: "", All: []string{"qa", "stg"}},
			args:            []string{"--env="},
			wantErr:         "please provide a value for 'env'",
		},
		{
			name:            "do nothing on input bool refusal",
			inputBoolResult: false,
			env:             envCompleted,
			envResultInFile: envCompleted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.envFileNil {
				jsonData, _ := json.Marshal(tt.env)
				err := ioutil.WriteFile(envFile, jsonData, os.ModePerm)
				assert.NoError(t, err)
			}

			listMock := &mocks.InputListMock{}
			boolMock := &mocks.InputBoolMock{}
			listMock.On("List", mock.Anything, mock.Anything, mock.Anything).Return(tt.inputListString, tt.inputListError)
			boolMock.On("Bool", mock.Anything, mock.Anything, mock.Anything).Return(tt.inputBoolResult, tt.inputBoolError)

			cmd := NewDeleteEnvCmd(envFindRemover, boolMock, listMock)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if err != nil {
				assert.Equal(t, err.Error(), tt.wantErr)
			} else {
				assert.Empty(t, tt.wantErr)

				assert.FileExists(t, envFile)

				envResult, err := envFinder.Find()
				assert.NoError(t, err)
				assert.Equal(t, tt.envResultInFile, envResult)
			}
		})
	}
}
