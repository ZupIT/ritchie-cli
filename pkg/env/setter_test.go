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

func TestSet(t *testing.T) {
	tmp := os.TempDir()
	ritHomeDir := filepath.Join(tmp, ".rit")
	file := stream.NewFileManager()
	finder := NewFinder(ritHomeDir, file)
	_ = os.MkdirAll(ritHomeDir, os.ModePerm)
	defer os.RemoveAll(ritHomeDir)

	type in struct {
		file      stream.FileWriter
		envFinder Finder
		env       string
	}

	type out struct {
		want Holder
		err  error
	}

	tests := []struct {
		name string
		in   in
		out  *out
	}{
		{
			name: "new dev env",
			in: in{
				file:      file,
				envFinder: finder,
				env:       dev,
			},
			out: &out{
				want: Holder{Current: dev, All: []string{dev}},
				err:  nil,
			},
		},
		{
			name: "no duplicate env",
			in: in{
				file:      file,
				envFinder: finder,
				env:       dev,
			},
			out: &out{
				want: Holder{Current: dev, All: []string{dev}},
				err:  nil,
			},
		},
		{
			name: "new qa env",
			in: in{
				file:      file,
				envFinder: finder,
				env:       qa,
			},
			out: &out{
				want: Holder{Current: qa, All: []string{dev, qa}},
				err:  nil,
			},
		},
		{
			name: "default env",
			in: in{
				file:      file,
				envFinder: finder,
				env:       Default,
			},
			out: &out{
				want: Holder{Current: "", All: []string{dev, qa}},
				err:  nil,
			},
		},
		{
			name: "find env error",
			in: in{
				file:      file,
				envFinder: findEnvMock{err: errors.New("error to find env")},
				env:       qa,
			},
			out: &out{
				err: errors.New("error to find env"),
			},
		},
		{
			name: "write env error",
			in: in{
				file: sMock.FileWriterCustomMock{
					WriteMock: func(path string, content []byte) error {
						return errors.New("error to write envs file")
					},
				},
				envFinder: finder,
				env:       qa,
			},
			out: &out{
				err: errors.New("error to write envs file"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in
			out := tt.out

			setter := NewSetter(ritHomeDir, in.envFinder, in.file)
			got, err := setter.Set(in.env)

			if out.err != nil && out.err.Error() != err.Error() {
				t.Errorf("Set(%s) got %v, want %v", tt.name, err, out.err)
			}
			if !reflect.DeepEqual(out.want, got) {
				t.Errorf("Set(%s) got %v, want %v", tt.name, got, out.want)
			}
		})
	}
}
