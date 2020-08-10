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

func TestSet(t *testing.T) {
	type out struct {
		want      MetricsHolder
		err       error
		waitError bool
	}
	err := errors.New("some error")

	tests := []struct {
		name       string
		in         string
		out        *out
		FileWriter stream.FileWriter
	}{
		{
			name: "Set enabled metrics",
			in:   "enabled",
			out: &out{
				want:      MetricsHolder{Current: "enabled"},
				err:       nil,
				waitError: false,
			},
			FileWriter: sMocks.FileWriterCustomMock{
				WriteMock: func(path string, content []byte) error {
					return nil
				},
			},
		},
		{
			name: "Set disabled Metrics",
			in:   "disabled",
			out: &out{
				want:      MetricsHolder{Current: "disabled"},
				err:       nil,
				waitError: false,
			},
			FileWriter: sMocks.FileWriterCustomMock{
				WriteMock: func(path string, content []byte) error {
					return nil
				},
			},
		},
		{
			name: "Error writing the Metrics file",
			in:   DefaultMetrics,
			out: &out{
				want:      MetricsHolder{Current: DefaultMetrics},
				err:       err,
				waitError: true,
			},
			FileWriter: sMocks.FileWriterCustomMock{
				WriteMock: func(path string, content []byte) error {
					return err
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setter := NewSetter("any/path", tt.FileWriter)

			in := tt.in
			out := tt.out

			got, err := setter.Set(in)
			if err != nil && !tt.out.waitError {
				t.Errorf("Set(%s) - Execution error - got %v, want %v", tt.name, err, out.err)
			}
			if !reflect.DeepEqual(out.want, got) {
				t.Errorf("Set(%s) - Error in the expected response -  got %v, want %v", tt.name, got, out.want)
			}
		})
	}

}
