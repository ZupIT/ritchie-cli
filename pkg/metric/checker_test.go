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

package metric

import (
	"errors"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	sMocks "github.com/ZupIT/ritchie-cli/pkg/stream/mocks"
)

func Test_Check(t *testing.T) {
	type in struct {
		file stream.FileWriteReadExister
		prompt prompt.InputList
	}

	var tests = []struct {
		name           string
		expectedResult bool
		in             in
	}{
		{
			name:           "success case expecting true",
			expectedResult: true,
			in: in{
				file: sMocks.FileWriteReadExisterCustomMock{
					ReadMock: func(path string) ([]byte, error) {
						return []byte("yes"), nil
					},
					ExistsMock: func(path string) bool {
						return true
					},
				},
			},
		},
		{
			name:           "success case expecting false",
			expectedResult: false,
			in: in{
				file: sMocks.FileWriteReadExisterCustomMock{
					ReadMock: func(path string) ([]byte, error) {
						return []byte("no"), nil
					},
					ExistsMock: func(path string) bool {
						return true
					},
				},
			},
		},
		{
			name:           "success case when metrics file doesn't exist",
			expectedResult: false,
			in: in{
				file: sMocks.FileWriteReadExisterCustomMock{
					ExistsMock: func(path string) bool {
						return false
					},
					WriteMock: func(path string, content []byte) error {
						return nil
					},
				},
			},
		},
		{
			name:           "success case expecting false when error reading file",
			expectedResult: false,
			in: in{
				file: sMocks.FileWriteReadExisterCustomMock{
					ReadMock: func(path string) ([]byte, error) {
						return nil, errors.New("error reading file")
					},
					ExistsMock: func(path string) bool {
						return true
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checker := NewChecker(tt.in.file)
			result := checker.Check()

			if result != tt.expectedResult {
				t.Errorf("behavior test failed: %s\nwant: %t | got: %t", tt.name, tt.expectedResult, result)
			}
		})
	}
}

