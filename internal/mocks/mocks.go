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
	"io"

	"github.com/stretchr/testify/mock"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/git"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
)

type ContextFinderMock struct {
	mock.Mock
}

func (cf *ContextFinderMock) Find() (rcontext.ContextHolder, error) {
	args := cf.Called()

	return args.Get(0).(rcontext.ContextHolder), args.Error(1)
}

type DetailManagerMock struct {
	mock.Mock
}

func (d *DetailManagerMock) LatestTag(repo formula.Repo) string {
	args := d.Called(repo)

	return args.String(0)
}

type RepositoriesMock struct {
	mock.Mock
}

func (m *RepositoriesMock) Zipball(info git.RepoInfo, version string) (io.ReadCloser, error) {
	args := m.Called(info, version)

	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *RepositoriesMock) Tags(info git.RepoInfo) (git.Tags, error) {
	args := m.Called(info)

	return args.Get(0).(git.Tags), args.Error(1)
}

func (m *RepositoriesMock) LatestTag(info git.RepoInfo) (git.Tag, error) {
	args := m.Called(info)

	return args.Get(0).(git.Tag), args.Error(1)
}
