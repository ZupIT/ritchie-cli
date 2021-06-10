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
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
	sMocks "github.com/ZupIT/ritchie-cli/pkg/stream/mocks"
	"github.com/stretchr/testify/assert"
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
			name:    "return empty repo when fails on read",
			wantErr: false,
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

func TestMetricId(t *testing.T) {

	tests := []struct {
		name string
		in   []string
		out  string
	}{
		{
			name: "success filter input flag with docker flag",
			in:   []string{"cmd", "test", "login", "--username=dennis", "--password=123456", "--docker"},
			out:  "rit_test_login",
		},
		{
			name: "success filter input flag for credential",
			in:   []string{"cmd", "set", "credential", "--provider=github", "--fields=username,token", "--values=\"$USERNAME_CREDENTIAL\",\"$GITHUB_TOKEN\""},
			out:  "rit_set_credential",
		},
		{
			name: "success filter input flag for core command",
			in:   []string{"cmd", "add", "repo", "--provider=\"Github\"", "--name=\"formulas-insights\"", "--repoUrl=\"https://github.com/ZupIT/ritchie-formulas\"", "--priority=1"},
			out:  "rit_add_repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldArgs := os.Args
			os.Args = tt.in
			defer func() { os.Args = oldArgs }()

			got := metricId()
			if got != tt.out {
				t.Errorf("Unexpected return: got %v want %v", got, tt.out)
			}

			assert.Equal(t, got, tt.out)
		})
	}
}

func TestFlags(t *testing.T) {

	tests := []struct {
		name string
		in   []string
		out  []string
	}{
		{
			name: "success filter input flag",
			in:   []string{"cmd", "create", "formula"},
			out:  nil,
			//out:  []string{"username", "password", "docker"},
		},
		{
			name: "success filter input flag with docker flag",
			in:   []string{"cmd", "test", "login", "--username=dennis", "--password=123456", "--docker"},
			out:  []string{"docker"},
			//out:  []string{"username", "password", "docker"},
		},
		{
			name: "success filter input flag for credential",
			in:   []string{"cmd", "set", "credential", "--provider=github", "--fields=username,token", "--values=\"$USERNAME_CREDENTIAL\",\"$GITHUB_TOKEN\""},
			out:  nil,
			//out:  []string{"provider", "fields", "values"},
		},
		{
			name: "success filter input flag for core command",
			in:   []string{"cmd", "add", "repo", "--provider=\"Github\"", "--name=\"formulas-insights\"", "--repoUrl=\"https://github.com/ZupIT/ritchie-formulas\"", "--priority=1"},
			out:  nil,
			//out:  []string{"provider", "name", "repoUrl", "priority"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldArgs := os.Args
			os.Args = tt.in
			defer func() { os.Args = oldArgs }()

			got := flags()

			assert.Equal(t, got, tt.out)
		})
	}
}
