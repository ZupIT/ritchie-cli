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

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type FindManager struct {
	filePath string
	file     stream.FileReadExister
}

func NewFinder(homePath string, file stream.FileReadExister) FindManager {
	return FindManager{
		filePath: filepath.Join(homePath, FileName),
		file:     file,
	}
}

func (f FindManager) Find() (Holder, error) {
	envHolder := Holder{}

	if !f.file.Exists(f.filePath) {
		return envHolder, nil
	}

	b, err := f.file.Read(f.filePath)
	if err != nil {
		return envHolder, err
	}

	if err := json.Unmarshal(b, &envHolder); err != nil {
		return envHolder, err
	}

	return envHolder, nil
}
