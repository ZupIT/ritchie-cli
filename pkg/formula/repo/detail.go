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

package repo

import (
	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

type DetailManager struct {
	repoProviders formula.RepoProviders
}

func NewDetail(repoProviders formula.RepoProviders) DetailManager {
	return DetailManager{repoProviders}
}

func (dm DetailManager) LatestTag(repo formula.Repo) string {
	formulaGit := dm.repoProviders.Resolve(repo.Provider)

	repoInfo := formulaGit.NewRepoInfo(repo.URL, repo.Token)
	tag, err := formulaGit.Repos.LatestTag(repoInfo)
	if err != nil {
		return ""
	}

	return tag.Name
}
