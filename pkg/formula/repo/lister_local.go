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
	"os"
	"path/filepath"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

type ListerLocalManager struct {
	ritHome string
}

func NewListerLocal(ritHome string) ListerLocalManager {
	return ListerLocalManager{ritHome: ritHome}
}

// List method returns an empty formula.RepoName if there is no local folder on li.ritHome
func (li ListerLocalManager) List() (formula.RepoName, error) {
	localReposPath := filepath.Join(li.ritHome, reposDirName, "local")

	if _, err := os.Stat(localReposPath); os.IsNotExist(err) {
		return "", err
	}

	return "local", nil
}
