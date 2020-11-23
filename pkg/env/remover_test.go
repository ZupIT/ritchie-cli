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
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
	sMock "github.com/ZupIT/ritchie-cli/pkg/stream/mocks"
)

func TestRemove(t *testing.T) {
	tmp := os.TempDir()
	file := stream.NewFileManager()
	finder := NewFinder(tmp, file)
	setter := NewSetter(tmp, finder, file)

	_, err := setter.Set(dev)
	if err != nil {
		fmt.Sprintln("Error in Set")
		return
	}
	_, err = setter.Set(qa)
	if err != nil {
		fmt.Sprintln("Error in Set")
		return
	}

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
			name: "dev env",
			in: in{
				file:      file,
				envFinder: finder,
				env:       dev,
			},
			out: &out{
				want: Holder{Current: qa, All: []string{qa}},
				err:  nil,
			},
		},
		{
			name: "current env",
			in: in{
				file:      file,
				envFinder: finder,
				env:       CurrentEnv + qa,
			},
			out: &out{
				want: Holder{All: []string{}},
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

			remover := NewRemover(tmp, in.envFinder, in.file)
			got, err := remover.Remove(in.env)
			if out.err != nil && out.err.Error() != err.Error() {
				t.Errorf("Remove(%s) got %v, want %v", tt.name, err, out.err)
			}
			if !reflect.DeepEqual(out.want, got) {
				t.Errorf("Remove(%s) got %v, want %v", tt.name, got, out.want)
			}
		})
	}
}
