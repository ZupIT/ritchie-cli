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

	"github.com/ZupIT/ritchie-cli/pkg/slice/sliceutil"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type SetterManager struct {
	filePath string
	env      Finder
	file     stream.FileWriter
}

func NewSetter(homePath string, env Finder, file stream.FileWriter) Setter {
	return SetterManager{
		filePath: filepath.Join(homePath, FileName),
		env:      env,
		file:     file,
	}
}

func (s SetterManager) Set(env string) (Holder, error) {
	envHolder, err := s.env.Find()
	if err != nil {
		return Holder{}, err
	}

	envHolder.Current = strings.ReplaceAll(env, Default, "")
	if env != Default {
		if envHolder.All == nil {
			envHolder.All = make([]string, 0)
		}

		if !sliceutil.Contains(envHolder.All, env) {
			envHolder.All = append(envHolder.All, env)
		}
	}

	b, err := json.Marshal(&envHolder)
	if err != nil {
		return Holder{}, err
	}
	if err := s.file.Write(s.filePath, b); err != nil {
		return Holder{}, err
	}

	return envHolder, nil
}
