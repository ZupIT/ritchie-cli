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

func TestNewSetEnvCmd(t *testing.T) {
	tmpDir := os.TempDir()
	ritHomeDir := filepath.Join(tmpDir, ".rit")
	envFile := filepath.Join(ritHomeDir, env.FileName)
	_ = os.MkdirAll(ritHomeDir, os.ModePerm)
	defer os.RemoveAll(ritHomeDir)

	fileManager := stream.NewFileManager()

	envFinder := env.NewFinder(ritHomeDir, fileManager)
	envSetter := env.NewSetter(ritHomeDir, envFinder, fileManager)
	envFindSetter := env.NewFindSetter(envFinder, envSetter)

	envExistingEnv := env.Holder{Current: "", All: []string{"existingEnv"}}

	type in struct {
		args            []string
		envList         env.Holder
		envResult       env.Holder
		inputListString string
		inputListErr    error
		inputText       string
		inputTextErr    error
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "success prompt new environment",
			in: in{
				args:            []string{},
				envResult:       env.Holder{Current: "newEnv", All: []string{"newEnv"}},
				inputListString: newEnv,
				inputText:       "newEnv",
			},
		},
		{
			name: "success prompt existing environment",
			in: in{
				args:            []string{},
				envList:         envExistingEnv,
				envResult:       env.Holder{Current: "existingEnv", All: []string{"existingEnv"}},
				inputListString: "existingEnv",
			},
		},
		{
			name: "success flag new environment",
			in: in{
				args:      []string{"--env=newEnv"},
				envResult: env.Holder{Current: "newEnv", All: []string{"newEnv"}},
			},
		},
		{
			name: "success flag existing environment",
			in: in{
				args:      []string{"--env=existingEnv"},
				envList:   envExistingEnv,
				envResult: env.Holder{Current: "existingEnv", All: []string{"existingEnv"}},
			},
		},
		{
			name: "error to list env",
			in: in{
				args:         []string{},
				inputListErr: errors.New("error to list env"),
			},
			want: errors.New("error to list env"),
		},
		{
			name: "error to input text",
			in: in{
				args:            []string{},
				inputListString: newEnv,
				inputTextErr:    errors.New("error to input text"),
			},
			want: errors.New("error to input text"),
		},
		{
			name: "error on empty flag",
			in: in{
				args: []string{"--env="},
			},
			want: errors.New("please provide a value for 'env'"),
		},
		{
			name: "error to resolve flag",
			in: in{
				args: []string{"--wrongFlag=newEnv"},
			},
			want: errors.New("unknown flag: --wrongFlag"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tt.in.envList)
			err := ioutil.WriteFile(envFile, jsonData, os.ModePerm)
			assert.NoError(t, err)

			inputListMock := new(mocks.InputListMock)
			inputListMock.On("List", mock.Anything, mock.Anything, mock.Anything).Return(tt.in.inputListString, tt.in.inputListErr)
			inputTextMock := new(mocks.InputTextMock)
			inputTextMock.On("Text", mock.Anything, mock.Anything, mock.Anything).Return(tt.in.inputText, tt.in.inputTextErr)

			cmd := NewSetEnvCmd(envFindSetter, inputTextMock, inputListMock)
			// TODO: remove stdin flag after  deprecation
			cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
			cmd.SetArgs(tt.in.args)

			got := cmd.Execute()
			if got != nil {
				assert.Equal(t, tt.want, got)
			} else {
				assert.Empty(t, tt.want)

				envResult, err := envFinder.Find()
				assert.NoError(t, err)
				assert.Equal(t, tt.in.envResult, envResult)
			}
		})
	}
}
