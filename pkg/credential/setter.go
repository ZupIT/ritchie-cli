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

package credential

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type SetManager struct {
	homePath string
	env      env.Finder
	dir      stream.DirCreater
	file     stream.FileWriter
}

func NewSetter(
	homePath string,
	env env.Finder,
	dir stream.DirCreater,
	file stream.FileWriter,
) SetManager {
	return SetManager{
		homePath: homePath,
		env:      env,
		dir:      dir,
		file:     file,
	}
}

func (s SetManager) Set(cred Detail) error {
	envHolder, err := s.env.Find()
	if err != nil {
		return err
	}
	if envHolder.Current == "" {
		envHolder.Current = env.Default
	}

	cb, err := json.Marshal(cred)
	if err != nil {
		return err
	}

	dir := filepath.Join(s.homePath, credentialDir, envHolder.Current)
	if err := s.dir.Create(dir); err != nil {
		return err
	}

	credFile := filepath.Join(s.homePath, credentialDir, envHolder.Current, cred.Service)
	if err := ioutil.WriteFile(credFile, cb, os.ModePerm); err != nil {
		return err
	}

	return nil

}
