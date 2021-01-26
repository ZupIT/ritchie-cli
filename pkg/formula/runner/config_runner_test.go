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
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
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

	type in struct {
		ritHome string
		file    stream.FileWriteReadExister
	}

	type out struct {
		runType formula.RunnerType
		err     error
	}

	tests := []struct {
		name    string
		ritHome string
		in      in
		out     out
	}{
		{
			name:    "find config success",
			ritHome: tmpDir,
			in: in{
				ritHome: tmpDir,
				file:    fileManagerMock{rBytes: []byte("0"), exist: true},
			},
			out: out{
				runType: formula.LocalRun,
				err:     nil,
			},
		},
		{
			name: "find config not found error",
			in: in{
				ritHome: tmpDir,
				file:    fileManagerMock{exist: false},
			},
			out: out{
				runType: formula.DefaultRun,
				err:     ErrConfigNotFound,
			},
		},
		{
			name: "find config read error",
			in: in{
				ritHome: tmpDir,
				file:    fileManagerMock{rErr: errors.New("read config error"), exist: true},
			},
			out: out{
				runType: formula.DefaultRun,
				err:     errors.New("read config error"),
			},
		},
		{
			name: "find config invalid runType",
			in: in{
				ritHome: tmpDir,
				file:    fileManagerMock{rBytes: []byte("error"), exist: true},
			},
			out: out{
				runType: formula.DefaultRun,
				err:     errors.New("strconv.Atoi: parsing \"error\": invalid syntax"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewConfigManager(tt.in.ritHome)
			got, err := config.Find()

			if (tt.out.err != nil && err == nil) || err != nil && err.Error() != tt.out.err.Error() {
				t.Errorf("Find(%s) got %v, want %v", tt.name, err, tt.out.err)
			}

			if !reflect.DeepEqual(tt.out.runType, got) {
				t.Errorf("Find(%s) got %v, want %v", tt.name, got, tt.out.runType)
			}
		})
	}
}
