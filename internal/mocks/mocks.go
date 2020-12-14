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

	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/git"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
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

type RepoListerAdderMock struct {
	mock.Mock
}

func (r *RepoListerAdderMock) List() (formula.Repos, error) {
	args := r.Called()

	return args.Get(0).(formula.Repos), args.Error(1)
}

func (r *RepoListerAdderMock) Add(repo formula.Repo) error {
	args := r.Called(repo)

	return args.Error(0)
}

type TutorialFindSetterMock struct {
	mock.Mock
}

func (t *TutorialFindSetterMock) Find() (rtutorial.TutorialHolder, error) {
	args := t.Called()

	return args.Get(0).(rtutorial.TutorialHolder), args.Error(1)
}

func (t *TutorialFindSetterMock) Set(tutorial string) (rtutorial.TutorialHolder, error) {
	args := t.Called(tutorial)

	return args.Get(0).(rtutorial.TutorialHolder), args.Error(1)
}

type GitRepositoryMock struct {
	mock.Mock
}

func (g *GitRepositoryMock) Zipball(info git.RepoInfo, version string) (io.ReadCloser, error) {
	args := g.Called(info, version)

	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (g *GitRepositoryMock) Tags(info git.RepoInfo) (git.Tags, error) {
	args := g.Called(info)

	return args.Get(0).(git.Tags), args.Error(1)
}

func (g *GitRepositoryMock) LatestTag(info git.RepoInfo) (git.Tag, error) {
	args := g.Called(info)

	return args.Get(0).(git.Tag), args.Error(1)
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

type InputTextValidatorMock struct {
	mock.Mock
}

func (i *InputTextValidatorMock) Text(name string, validate func(interface{}) error, helper ...string) (string, error) {
	args := i.Called(name, validate)

	return args.String(0), args.Error(1)
}
