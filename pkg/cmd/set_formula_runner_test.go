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
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ZupIT/ritchie-cli/internal/mocks"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/runner"
)

func TestSetFormulaRunnerCmd(t *testing.T) {
	tmpDir := os.TempDir()
	ritHome := filepath.Join(tmpDir, "runner")
	_ = os.Mkdir(ritHome, os.ModePerm)
	defer os.RemoveAll(ritHome)

	configManager := runner.NewConfigManager(ritHome)
	runnerFile := filepath.Join(ritHome, runner.FileName)

	var tests = []struct {
		name    string
		args    []string
		runner  string
		listErr error
		err     error
	}{
		{
			name:   "success prompt set formula run",
			args:   []string{},
			runner: formula.LocalRun.String(),
		},
		{
			name:    "error on list",
			args:    []string{},
			listErr: errors.New("list error"),
			err:     errors.New("list error"),
		},
		{
			name:   "success on flags",
			args:   []string{"--runner=local"},
			runner: formula.LocalRun.String(),
		},
		{
			name: "fail when missing flag",
			args: []string{"--runner="},
			err:  errors.New(missingFlagText(runnerFlagName)),
		},
		{
			name: "fail for wrong flag name",
			args: []string{"--runner=invalid"},
			err:  ErrInvalidRunType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Remove(runnerFile)
			inputList := &mocks.InputListMock{}
			inputList.On(
				"List", "Select a default formula run type", mock.Anything, mock.Anything,
			).Return(tt.runner, tt.listErr)

			cmd := NewSetFormulaRunnerCmd(configManager, inputList)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()

			if err != nil {
				assert.Equal(t, tt.err, err)
			} else {
				assert.Nil(t, tt.err)
				assert.FileExists(t, runnerFile)

				data, err := ioutil.ReadFile(runnerFile)
				assert.NoError(t, err)
				runType, err := strconv.Atoi(string(data))
				assert.NoError(t, err)
				assert.Equal(t, tt.runner, formula.RunnerType(runType).String())
			}
		})
	}
}
