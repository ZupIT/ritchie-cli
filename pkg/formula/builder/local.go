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

package builder

import (
	"path/filepath"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/repo/repoutil"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type Initializer interface {
	Init(workspaceDir string, repoName string) (string, error)
}

type LocalManager struct {
	ritHome string
	dir     stream.DirCreateListCopyRemover
	repo    formula.RepositoryAdder
}

func NewBuildLocal(
	ritHome string,
	dir stream.DirCreateListCopyRemover,
	repo formula.RepositoryAdder,
) LocalManager {
	return LocalManager{ritHome: ritHome, dir: dir, repo: repo}
}

func (m LocalManager) Build(info formula.BuildInfo) error {
	dest, err := m.Init(info.Workspace.Dir, info.Workspace.Name)
	if err != nil {
		return err
	}

	formulaSrc := strings.ReplaceAll(info.FormulaPath, info.Workspace.Dir, dest)
	formulaBin := filepath.Join(formulaSrc, "bin")
	if err := m.dir.Remove(formulaBin); err != nil {
		return err
	}

	return nil
}

func (m LocalManager) Init(workspaceDir string, repoName string) (string, error) {
	repoNameStandard := repoutil.LocalName(repoName)
	dest := filepath.Join(m.ritHome, "repos", repoNameStandard.String())

	repo := formula.Repo{
		Provider:      "Local",
		Name:          repoNameStandard,
		Version:       "0.0.0",
		Url:           "local repository",
		Priority:      0,
		IsLocal:       true,
		LatestVersion: "0.0.0",
	}

	if err := m.dir.Create(dest); err != nil {
		return "", err
	}

	if err := m.copyWorkSpace(workspaceDir, dest); err != nil {
		return "", err
	}

	if err := m.repo.Add(repo); err != nil {
		return "", err
	}

	return dest, nil
}

func (m LocalManager) copyWorkSpace(workspacePath string, dest string) error {
	return m.dir.Copy(workspacePath, dest)
}
