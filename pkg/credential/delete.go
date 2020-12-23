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
	"os"
	"path/filepath"

	renv "github.com/ZupIT/ritchie-cli/pkg/env"
)

type DeleteManager struct {
	homePath string
	env      renv.Finder
}

func NewCredDelete(homePath string, env renv.Finder) DeleteManager {
	return DeleteManager{
		homePath: homePath,
		env:      env,
	}
}

func (d DeleteManager) Delete(service string) error {
	env, err := d.env.Find()
	if err != nil {
		return err
	}

	if env.Current == "" {
		env.Current = renv.Default
	}

	path := filepath.Join(d.homePath, credentialDir, env.Current, service)
	if err := os.Remove(path); err != nil {
		return err
	}
	return nil
}
