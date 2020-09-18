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

	"github.com/ZupIT/ritchie-cli/pkg/stream"
	sMocks "github.com/ZupIT/ritchie-cli/pkg/stream/mocks"
)

func Test_Collector(t *testing.T) {
	repoJson := `[{
		"provider": "Github",
		"name": "",
		"version": "2.6.0",
		"url": "https://github.com/ZupIT/ritchie-formulas",
		"priority": 0
	}]
`
	type in struct {
		userIdGen UserIdGenerator
		file      stream.FileReader
	}

	var tests = []struct {
		name    string
		wantErr bool
		in
	}{
		{
			name:    "success case",
			wantErr: false,
			in: in{
				userIdGen: UserIdGeneratorMock{
					GenerateMock: func() (UserId, error) {
						return "", nil
					}},
				file: sMocks.FileReaderCustomMock{
					ReadMock: func(path string) ([]byte, error) {
						return []byte(repoJson), nil
					}},
			},
		},
		{
			name:    "fails when generator returns an error",
			wantErr: true,
			in: in{
				userIdGen: UserIdGeneratorMock{
					GenerateMock: func() (UserId, error) {
						return "", errors.New("error generating id")
					}},
			},
		},
		{
			name:    "fails on read file",
			wantErr: true,
			in: in{
				userIdGen: UserIdGeneratorMock{
					GenerateMock: func() (UserId, error) {
						return "", nil
					}},
				file: sMocks.FileReaderCustomMock{
					ReadMock: func(path string) ([]byte, error) {
						return nil, errors.New("error reading file")
					}},
			},
		},
		{
			name:    "fails on unmarshal json",
			wantErr: true,
			in: in{
				userIdGen: UserIdGeneratorMock{
					GenerateMock: func() (UserId, error) {
						return "", nil
					}},
				file: sMocks.FileReaderCustomMock{
					ReadMock: func(path string) ([]byte, error) {
						return []byte(`{"someRandomKey":"andTheRespectiveRandomValue"}`), nil
					}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := NewDataCollector(tt.in.userIdGen, "", tt.in.file)
			_, err := collector.Collect(1, "2.0.0")
			if (err != nil) != tt.wantErr {
				t.Errorf("execution test failed: %s\nwant error: %t | got: %s", tt.name, tt.wantErr, err)
			}
		})
	}

}

type UserIdGeneratorMock struct {
	GenerateMock func() (UserId, error)
}

func (us UserIdGeneratorMock) Generate() (UserId, error) {
	return us.GenerateMock()
}
