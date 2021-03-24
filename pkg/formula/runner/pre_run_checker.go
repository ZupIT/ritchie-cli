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

package runner

import (
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/formula/repo"
)

const ErrPreRunCheckerVersion = "Failed to run formula, this formula needs run in the last version of repository.\n\tCurrent version: %s\n\tLatest version: %s"

type PreRunCheckerManager struct {
	repoList repo.ListManager
}

func NewPreRunBuilderChecker(listManager repo.ListManager) PreRunCheckerManager {
	return PreRunCheckerManager{listManager}
}

func (pr *PreRunCheckerManager) CheckVersionCompliance(repoName string, requireLatestVersion bool) error {
	if requireLatestVersion {
		repos, _ := pr.repoList.List()
		repo, _ := repos.Get(repoName)
		if repo.Version.String() != repo.LatestVersion.String() {
			return fmt.Errorf(ErrPreRunCheckerVersion, repo.Version.String(), repo.LatestVersion.String())
		}
	}
	return nil
}
