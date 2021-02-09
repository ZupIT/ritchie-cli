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
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

func TestCreate(t *testing.T) {
	tmpDir := os.TempDir()
	ritHome := filepath.Join(tmpDir, "create")
	ritInvalidHome := filepath.Join(tmpDir, "invalid")
	_ = os.Mkdir(ritHome, os.ModePerm)
	defer os.RemoveAll(ritHome)

	tests := []struct {
		name    string
		ritHome string
		err     string
	}{
		{
			name:    "create config success",
			ritHome: ritHome,
		},
		{
			name:    "create config write error",
			ritHome: ritInvalidHome,
			err: fmt.Sprintf(
				"open %s: no such file or directory",
				filepath.Join(ritInvalidHome, FileName),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := filepath.Join(ritHome, FileName)
			_ = os.Remove(file)
			config := NewConfigManager(tt.ritHome)
			got := config.Create(formula.LocalRun)

			if got != nil {
				assert.EqualError(t, got, tt.err)
			} else {
				assert.Empty(t, tt.err)
				assert.FileExists(t, file)

				data, err := ioutil.ReadFile(file)
				assert.NoError(t, err)

				runType, err := strconv.Atoi(string(data))
				assert.NoError(t, err)
				assert.Equal(t, formula.LocalRun.Int(), runType)
			}
		})
	}
}

func TestFind(t *testing.T) {
	tmpDir := os.TempDir()
	ritHome := filepath.Join(tmpDir, "find")
	ritInvalidHome := filepath.Join(tmpDir, "invalid")
	_ = os.Mkdir(ritHome, os.ModePerm)
	defer os.RemoveAll(ritHome)

	tests := []struct {
		name        string
		ritHome     string
		fileContent string
		runner      formula.RunnerType
		err         string
	}{
		{
			name:        "find config success",
			ritHome:     ritHome,
			fileContent: "0",
			runner:      formula.LocalRun,
		},
		{
			name:    "fail finding file",
			ritHome: ritInvalidHome,
			err:     ErrConfigNotFound.Error(),
		},
		{
			name:        "fail invalid runType",
			ritHome:     ritHome,
			fileContent: "error",
			err:         "strconv.Atoi: parsing \"error\": invalid syntax",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := filepath.Join(tt.ritHome, FileName)
			_ = ioutil.WriteFile(file, []byte(tt.fileContent), os.ModePerm)

			config := NewConfigManager(tt.ritHome)
			got, err := config.Find()

			if err != nil {
				assert.EqualError(t, err, tt.err)
			} else {
				assert.Empty(t, tt.err)
				assert.Equal(t, got, tt.runner)
			}
		})
	}
}
