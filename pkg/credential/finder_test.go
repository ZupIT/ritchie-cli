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

package credential

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestFind(t *testing.T) {

	fileManager := stream.NewFileManager()
	tmp := os.TempDir()
	setter := NewSetter(tmp, ctxFinder)
	_ = setter.Set(githubCred)

	type out struct {
		cred Detail
		err  error
	}

	type in struct {
		homePath  string
		ctxFinder rcontext.Finder
		file      stream.FileReader
		provider  string
	}

	tests := []struct {
		name string
		in   in
		out  out
	}{
		{
			name: "Run with success",
			in: in{
				homePath:  tmp,
				ctxFinder: ctxFinder,
				file:      fileManager,
				provider:  githubCred.Service,
			},
			out: out{
				cred: githubCred,
				err:  nil,
			},
		},
		{
			name: "Return err when file not exist",
			in: in{
				homePath:  tmp,
				ctxFinder: ctxFinder,
				file:      fileManager,
				provider:  "aws",
			},
			out: out{
				cred: Detail{Credential: Credential{}},
				err:  errors.New(prompt.Red(fmt.Sprintf(errNotFoundTemplate, "aws"))),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := tt.out
			finder := NewFinder(tt.in.homePath, tt.in.ctxFinder, tt.in.file)
			got, err := finder.Find(tt.in.provider)
			if err != nil && err.Error() != out.err.Error() {
				t.Errorf("Find(%s) got %v, want %v", tt.name, err, out.err)
			}

			if !reflect.DeepEqual(out.cred, got) {
				t.Errorf("Find(%s) got %v, want %v", tt.name, got, out.cred)
			}
		})
	}
}
