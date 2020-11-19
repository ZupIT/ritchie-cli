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
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"os"
	"testing"

	stream "github.com/ZupIT/ritchie-cli/pkg/stream/mocks"
)

var (
	githubCred = Detail{Service: "github"}
	streamMock = stream.FileReadExisterCustomMock{
		ReadMock: func(path string) ([]byte, error) {
			return []byte("{\"current_context\":\"default\"}"), nil
		},
		ExistsMock: func(path string) bool {
			return true
		},
	}
	ctxFinder = env.FindManager{filePath: "", file: streamMock}
)

func TestSet(t *testing.T) {

	tmp := os.TempDir()
	setter := NewSetter(tmp, ctxFinder)
	tests := []struct {
		name string
		in   Detail
		out  error
	}{
		{
			name: "github credential",
			in:   githubCred,
			out:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := setter.Set(tt.in)
			if got != tt.out {
				t.Errorf("Set(%s) got %v, want %v", tt.name, got, tt.out)
			}
		})
	}
}
