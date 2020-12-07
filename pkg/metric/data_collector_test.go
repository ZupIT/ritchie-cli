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

	"github.com/stretchr/testify/assert"

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
					GenerateMock: func() UserID {
						return ""
					}},
				file: sMocks.FileReaderCustomMock{
					ReadMock: func(path string) ([]byte, error) {
						return []byte(repoJson), nil
					}},
			},
		},
		{
			name:    "return empty repo when fails on read",
			wantErr: false,
			in: in{
				userIdGen: UserIdGeneratorMock{
					GenerateMock: func() UserID {
						return ""
					}},
				file: sMocks.FileReaderCustomMock{
					ReadMock: func(path string) ([]byte, error) {
						return nil, errors.New("error reading file")
					}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := NewDataCollector(UserIdGeneratorMock{
				GenerateMock: func() UserID {
					return ""
				}}, "", tt.in.file)
			user := collector.CollectUserState("2.0.0")
			assert.NotEmpty(t, user)

			command := collector.CollectCommandData(1)
			assert.NotEmpty(t, command)
		})
	}
}

type UserIdGeneratorMock struct {
	GenerateMock func() UserID
}

func (us UserIdGeneratorMock) Generate() UserID {
	return us.GenerateMock()
}
