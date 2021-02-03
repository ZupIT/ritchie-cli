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

package env

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
	sMock "github.com/ZupIT/ritchie-cli/pkg/stream/mocks"
)

func TestFind(t *testing.T) {
	tmp := os.TempDir()
	ritHomeDir := filepath.Join(tmp, ".rit")
	_ = os.MkdirAll(ritHomeDir, os.ModePerm)
	defer os.RemoveAll(ritHomeDir)

	type in struct {
		holder          Holder
		FileReadExister stream.FileReadExister
	}

	type out struct {
		err  error
		want Holder
	}

	tests := []struct {
		name string
		in   *in
		out  *out
	}{
		{
			name: "default env and existing env file",
			in: &in{
				holder: Holder{Current: ""},
				FileReadExister: sMock.FileReadExisterCustomMock{
					ReadMock: func(path string) ([]byte, error) {
						return []byte("{\"current_env\":\"default\"}"), nil
					},
					ExistsMock: func(path string) bool {
						return true
					},
				},
			},
			out: &out{
				want: Holder{Current: "default"},
				err:  nil,
			},
		},
		{
			name: "default env and missing env file",
			in: &in{
				holder: Holder{Current: ""},
				FileReadExister: sMock.FileReadExisterCustomMock{
					ReadMock: func(path string) ([]byte, error) {
						return []byte("{\"current_env\":\"default\"}"), nil
					},
					ExistsMock: func(path string) bool {
						return false
					},
				},
			},
			out: &out{
				want: Holder{Current: ""},
			},
		},
		{
			name: "default env and error on read file",
			in: &in{
				holder: Holder{Current: ""},
				FileReadExister: sMock.FileReadExisterCustomMock{
					ReadMock: func(path string) ([]byte, error) {
						return []byte(""), errors.New("error reading file")
					},
					ExistsMock: func(path string) bool {
						return true
					},
				},
			},
			out: &out{
				want: Holder{Current: ""},
				err:  errors.New("error reading file"),
			},
		},
		{
			name: "default env and incorrect json",
			in: &in{
				holder: Holder{Current: "default"},
				FileReadExister: sMock.FileReadExisterCustomMock{
					ReadMock: func(path string) ([]byte, error) {
						return []byte(""), nil
					},
					ExistsMock: func(path string) bool {
						return true
					},
				},
			},
			out: &out{
				want: Holder{Current: ""},
				err:  errors.New("unexpected end of JSON input"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			finder := NewFinder(ritHomeDir, tt.in.FileReadExister)
			out := tt.out
			got, err := finder.Find()
			if out.err != nil && out.err.Error() != err.Error() {
				t.Errorf("Find(%s) - Execution error - got %v, want %v", tt.name, err, out.err)
			}
			if !reflect.DeepEqual(out.want, got) {
				t.Errorf("Find(%s) - Error in the expected response - got %v, want %v", tt.name, got, out.want)
			}
		})
	}
}

type findEnvMock struct {
	data Holder
	err  error
}

func (f findEnvMock) Find() (Holder, error) {
	return f.data, f.err
}
