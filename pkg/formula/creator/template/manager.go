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

package template

import (
	"sort"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

var _ Manager = DefaultManager{}

type Manager interface {
	Languages() ([]string, error)
}

func NewManager(ritchieHome string, dir stream.DirChecker) Manager {
	return DefaultManager{ritchieHome, dir}
}

type DefaultManager struct {
	ritchieHome string
	dir         stream.DirChecker
}

func (tm DefaultManager) Languages() ([]string, error) {
	languages, err := AssetDir("templates/languages")
	if err != nil {
		return nil, err
	}

	sort.Strings(languages)

	return languages, nil
}
