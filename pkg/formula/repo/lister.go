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
	"encoding/json"
	"path/filepath"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type ListManager struct {
	ritHome string
	file    stream.FileReadExister
}

func NewLister(ritHome string, file stream.FileReadExister) ListManager {
	return ListManager{ritHome: ritHome, file: file}
}

// List method returns an empty formula.Repos if there is no repositories.json
func (li ListManager) List() (formula.Repos, error) {
	repos := formula.Repos{}
	reposFilePath := filepath.Join(li.ritHome, reposDirName, reposFileName)
	if !li.file.Exists(reposFilePath) {
		return repos, nil
	}

	file, err := li.file.Read(reposFilePath)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(file, &repos); err != nil {
		return nil, err
	}

	return repos, nil
}
