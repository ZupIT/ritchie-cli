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
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/stretchr/testify/mock"
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

type RepoListerAdderMock struct {
	mock.Mock
}

func (r *RepoListerAdderMock) List() (formula.Repos, error) {
	args := r.Called()

	return args.Get(0).(formula.Repos), args.Error(1)
}

func (r *RepoListerAdderMock) Add(repo formula.Repo) error {
	args := r.Called(repo)

	return args.Error(1)
}

type InputURLMock struct {
	mock.Mock
}

func (i *InputURLMock) URL(name, defaultValue string) (string, error) {
	args := i.Called(name, defaultValue)

	return args.String(0), args.Error(1)
}

type InputBoolMock struct {
	mock.Mock
}

func (i *InputBoolMock) Bool(name string, items []string, helper ...string) (bool, error) {
	args := i.Called(name, items, helper)

	return args.Bool(0), args.Error(1)
}

type InputListMock struct {
	mock.Mock
}

func (i *InputListMock) List(name string, items []string, helper ...string) (string, error) {
	args := i.Called(name, items, helper)

	return args.String(0), args.Error(1)
}

type InputIntMock struct {
	mock.Mock
}

func (i *InputIntMock) Int(name string, helper ...string) (int64, error) {
	args := i.Called(name, helper)

	return args.Get(0).(int64), args.Error(1)
}

type InputPasswordMock struct {
	mock.Mock
}

func (i *InputPasswordMock) Password(label string, helper ...string) (string, error) {
	args := i.Called(label, helper)

	return args.String(0), args.Error(1)
}
