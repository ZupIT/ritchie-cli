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

	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

type UpdaterMock struct {
	mock.Mock
}

func (u *UpdaterMock) List() (formula.Repos, error) {
	args := u.Called()

	return args.Get(0).(formula.Repos), args.Error(1)
}

func (u *UpdaterMock) Update(name formula.RepoName, version formula.RepoVersion, url string) error {
	args := u.Called(name, version, url)

	return args.Error(0)
}
