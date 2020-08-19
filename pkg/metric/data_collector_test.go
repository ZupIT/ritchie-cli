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
)

func Test_Collector(t *testing.T) {
	type in struct {
		userIdGen UserIdGenerator
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := NewDataCollector(tt.in.userIdGen)
			_, err := collector.Collect("2.0.0")
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
