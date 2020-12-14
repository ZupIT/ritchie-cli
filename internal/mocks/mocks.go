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

package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

type EnvFinderMock struct {
	mock.Mock
}

func (e *EnvFinderMock) Find() (env.Holder, error) {
	args := e.Called()

	return args.Get(0).(env.Holder), args.Error(1)
}

type DetailManagerMock struct {
	mock.Mock
}

func (d *DetailManagerMock) LatestTag(repo formula.Repo) string {
	args := d.Called(repo)

	return args.String(0)
}

type InputListMock struct {
	mock.Mock
}

func (m *InputListMock) List(string, []string, ...string) (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

type InputBoolMock struct {
	mock.Mock
}

func (m *InputBoolMock) Bool(string, []string, ...string) (bool, error) {
	args := m.Called()
	return args.Bool(0), nil
}
