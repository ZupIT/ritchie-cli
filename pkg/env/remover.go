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

package env

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type RemoveManager struct {
	filePath string
	env      Finder
	file     stream.FileWriter
}

func NewRemover(homePath string, env Finder, file stream.FileWriter) RemoveManager {
	return RemoveManager{
		filePath: filepath.Join(homePath, FileName),
		env:      env,
		file:     file,
	}
}

func (r RemoveManager) Remove(env string) (Holder, error) {
	envHolder, err := r.env.Find()
	if err != nil {
		return Holder{}, err
	}

	env = strings.ReplaceAll(env, Current, "")
	if envHolder.Current == env {
		envHolder.Current = ""
	}

	for i, e := range envHolder.All {
		if e == env {
			envHolder.All = append(envHolder.All[:i], envHolder.All[i+1:]...)
			break
		}
	}

	b, err := json.Marshal(&envHolder)
	if err != nil {
		return Holder{}, err
	}

	if err := r.file.Write(r.filePath, b); err != nil {
		return Holder{}, err
	}

	return envHolder, nil
}
