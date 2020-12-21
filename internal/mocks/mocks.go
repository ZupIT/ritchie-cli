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

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/template"
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

type InputTextMock struct {
	mock.Mock
}

func (i *InputTextMock) Text(name string, required bool, helper ...string) (string, error) {
	args := i.Called(name, required, helper)
	return args.String(0), args.Error(1)
}

type InputTextValidatorMock struct {
	mock.Mock
}

func (i *InputTextValidatorMock) Text(name string, validate func(interface{}) error, helper ...string) (string, error) {
	args := i.Called(name, validate)

	return args.String(0), args.Error(1)
}

type FormCreator struct {
	mock.Mock
}

func (f *FormCreator) Create(cf formula.Create) error {
	args := f.Called(cf)
	return args.Error(0)
}

func (f *FormCreator) Build(info formula.BuildInfo) error {
	args := f.Called(info)
	return args.Error(0)
}

type WorkspaceForm struct {
	mock.Mock
}

func (w *WorkspaceForm) Add(workspace formula.Workspace) error {
	args := w.Called(workspace)
	return args.Error(0)
}

func (w *WorkspaceForm) Delete(workspace formula.Workspace) error {
	args := w.Called(workspace)
	return args.Error(0)
}

func (w *WorkspaceForm) List() (formula.Workspaces, error) {
	args := w.Called()
	return args.Get(0).(formula.Workspaces), args.Error(1)
}

func (w *WorkspaceForm) Validate(workspace formula.Workspace) error {
	args := w.Called(workspace)
	return args.Error(0)
}

func (w *WorkspaceForm) CurrentHash(formulaPath string) (string, error) {
	args := w.Called(formulaPath)
	return args.String(0), args.Error(1)
}

func (w *WorkspaceForm) PreviousHash(formulaPath string) (string, error) {
	args := w.Called(formulaPath)
	return args.String(0), args.Error(1)
}

func (w *WorkspaceForm) UpdateHash(formulaPath string, hash string) error {
	args := w.Called(formulaPath, hash)
	return args.Error(0)
}

type RepoManager struct {
	mock.Mock
}

func (r *RepoManager) Add(repo formula.Repo) error {
	args := r.Called(repo)
	return args.Error(0)
}

func (r *RepoManager) List() (formula.Repos, error) {
	args := r.Called()
	return args.Get(0).(formula.Repos), args.Error(1)
}

func (r *RepoManager) Update(name formula.RepoName, version formula.RepoVersion) error {
	args := r.Called(name, version)
	return args.Error(0)
}

func (r *RepoManager) Delete(name formula.RepoName) error {
	args := r.Called(name)
	return args.Error(0)
}

func (r *RepoManager) SetPriority(name formula.RepoName, priority int) error {
	args := r.Called(name, priority)
	return args.Error(0)
}

func (r *RepoManager) Create(repo formula.Repo) error {
	args := r.Called(repo)
	return args.Error(0)
}

func (r *RepoManager) Write(repos formula.Repos) error {
	args := r.Called(repos)
	return args.Error(0)
}

func (r *RepoManager) LatestTag(repo formula.Repo) string {
	args := r.Called(repo)
	return args.String(0)
}

type FileManager struct {
	mock.Mock
}

func (f *FileManager) Exists(path string) bool {
	args := f.Called(path)
	return args.Bool(0)
}

func (f *FileManager) Read(path string) ([]byte, error) {
	args := f.Called(path)
	return args.Get(0).([]byte), args.Error(1)
}

func (f *FileManager) Write(path string, content []byte) error {
	args := f.Called(path, content)
	return args.Error(0)
}

func (f *FileManager) Create(path string, data io.ReadCloser) error {
	args := f.Called(path, data)
	return args.Error(0)
}

func (f *FileManager) Remove(path string) error {
	args := f.Called(path)
	return args.Error(0)
}

func (f *FileManager) List(file string) ([]string, error) {
	args := f.Called(file)
	return args.Get(0).([]string), args.Error(1)
}

func (f *FileManager) ListNews(oldPath, newPath string) ([]string, error) {
	args := f.Called(oldPath, newPath)
	return args.Get(0).([]string), args.Error(1)
}

func (f *FileManager) Copy(src, dst string) error {
	args := f.Called(src, dst)
	return args.Error(0)
}

func (f *FileManager) Append(path string, content []byte) error {
	args := f.Called(path, content)
	return args.Error(0)
}

func (f *FileManager) Move(oldPath, newPath string, files []string) error {
	args := f.Called(oldPath, newPath, files)
	return args.Error(0)
}

type TreeManager struct {
	mock.Mock
}

func (t *TreeManager) Tree() (map[formula.RepoName]formula.Tree, error) {
	args := t.Called()
	return args.Get(0).(map[formula.RepoName]formula.Tree), args.Error(1)
}

func (t *TreeManager) MergedTree(core bool) formula.Tree {
	args := t.Called(core)
	return args.Get(0).(formula.Tree)
}

func (t *TreeManager) Generate(repoPath string) (formula.Tree, error) {
	args := t.Called(repoPath)
	return args.Get(0).(formula.Tree), args.Error(1)
}

func (t *TreeManager) Check() []api.CommandID {
	args := t.Called()
	return args.Get(0).([]api.CommandID)
}

type TemplateManagerMock struct {
	mock.Mock
}

func (tm *TemplateManagerMock) Languages() ([]string, error) {
	args := tm.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (tm *TemplateManagerMock) LangTemplateFiles(lang string) ([]template.File, error) {
	args := tm.Called(lang)
	return args.Get(0).([]template.File), args.Error(1)
}

func (tm *TemplateManagerMock) ResolverNewPath(oldPath, newDir, lang, workspacePath string) (string, error) {
	args := tm.Called(oldPath, newDir, lang, workspacePath)
	return args.String(0), args.Error(1)
}

func (tm *TemplateManagerMock) Validate() error {
	args := tm.Called()
	return args.Error(0)
}
