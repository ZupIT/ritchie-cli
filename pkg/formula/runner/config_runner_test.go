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
	"os"
	"reflect"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestCreate(t *testing.T) {
	tmpDir := os.TempDir()

	type in struct {
		ritHome string
		file    stream.FileWriteReadExister
		runType formula.RunnerType
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "create config success",
			in: in{
				ritHome: tmpDir,
				file:    fileManagerMock{},
				runType: formula.Local,
			},
			want: nil,
		},
		{
			name: "create config write error",
			in: in{
				ritHome: tmpDir,
				file:    fileManagerMock{wErr: errors.New("error to write file")},
				runType: formula.Local,
			},
			want: errors.New("error to write file"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewConfigManager(tt.in.ritHome, tt.in.file)
			got := config.Create(tt.in.runType)

			if (tt.want != nil && got == nil) || got != nil && got.Error() != tt.want.Error() {
				t.Errorf("Create(%s) got %v, want %v", tt.name, got, tt.want)
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
		name string
		in   in
		out  out
	}{
		{
			name: "find config success",
			in: in{
				ritHome: tmpDir,
				file:    fileManagerMock{rBytes: []byte("0"), exist: true},
			},
			out: out{
				runType: formula.Local,
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
				runType: -1,
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
				runType: -1,
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
				runType: -1,
				err:     errors.New("strconv.Atoi: parsing \"error\": invalid syntax"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewConfigManager(tt.in.ritHome, tt.in.file)
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
