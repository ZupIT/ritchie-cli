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

package metric

import (
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type CheckManager struct {
	file stream.FileReadExister
}

func NewChecker(file stream.FileReadExister) CheckManager {
	return CheckManager{file: file}
}

func (c CheckManager) Check() (bool, error) {
	if !c.file.Exists(FilePath) {
		return false, nil
	}

	bytes, err := c.file.Read(FilePath)
	if err != nil {
		return false, err
	}

	result := true
	if string(bytes) == "no" {
		result = false
	}

	return result, nil
}

