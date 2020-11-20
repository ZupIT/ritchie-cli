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

	"github.com/ZupIT/ritchie-cli/pkg/env"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
)

type SetManager struct {
	homePath string
	ctx      env.Finder
}

func NewSetter(homePath string, cf env.Finder) SetManager {
	return SetManager{
		homePath: homePath,
		ctx:      cf,
	}
}

func (s SetManager) Set(cred Detail) error {
	ctx, err := s.ctx.Find()
	if err != nil {
		return err
	}
	if ctx.Current == "" {
		ctx.Current = env.DefaultEnv
	}

	cb, err := json.Marshal(cred)
	if err != nil {
		return err
	}

	dir := Dir(s.homePath, ctx.Current)
	if err := fileutil.CreateDirIfNotExists(dir, 0700); err != nil {
		return err
	}

	credFile := File(s.homePath, ctx.Current, cred.Service)
	if err := fileutil.WriteFilePerm(credFile, cb, 0600); err != nil {
		return err
	}

	return nil

}
