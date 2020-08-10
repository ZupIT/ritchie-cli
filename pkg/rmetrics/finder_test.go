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

package rmetrics

import (
	"errors"
	"reflect"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
	sMocks "github.com/ZupIT/ritchie-cli/pkg/stream/mocks"
)

var errReadingFile = errors.New("error reading file")

func TestFind(t *testing.T) {
	type out struct {
		err       error
		want      MetricsHolder
		wantError bool
	}

	tests := []struct {
		name string
		in   stream.FileReadExister
		out  *out
	}{
		{
			name: "With no metrics file",
			in: sMocks.FileReadExisterCustomMock{
				ExistsMock: func(path string) bool {
					return false
				},
			},
			out: &out{
				want:      MetricsHolder{Current: DefaultMetrics},
				err:       nil,
				wantError: false,
			},
		},
		{
			name: "With existing metrics file",
			in: sMocks.FileReadExisterCustomMock{
				ReadMock: func(path string) ([]byte, error) {
					return []byte("{\"metrics\":\"enabled\"}"), nil
				},
				ExistsMock: func(path string) bool {
					return true
				},
			},
			out: &out{
				want:      MetricsHolder{Current: "enabled"},
				err:       nil,
				wantError: false,
			},
		},
		{
			name: "Error reading the metrics file",
			in: sMocks.FileReadExisterCustomMock{
				ReadMock: func(path string) ([]byte, error) {
					return []byte(""), errReadingFile
				},
				ExistsMock: func(path string) bool {
					return true
				},
			},
			out: &out{
				want:      MetricsHolder{Current: DefaultMetrics},
				err:       errReadingFile,
				wantError: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			finder := NewFinder("any/path", tt.in)

			out := tt.out
			got, err := finder.Find()
			if err != nil && !tt.out.wantError {
				t.Errorf("%s - Execution error - got %v, want %v", tt.name, err, out.err)
			}
			if !reflect.DeepEqual(out.want, got) {
				t.Errorf("%s - Error in the expected response -  got %v, want %v", tt.name, got, out.want)
			}
		})
	}
}
