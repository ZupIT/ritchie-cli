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
	"errors"
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/env"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const errNotFoundTemplate = `
Fail to resolve credential for provider %s.
Try again after use:
	∙ rit set credential
`

type Finder struct {
	homePath  string
	ctxFinder env.Finder
	file      stream.FileReader
}

func NewFinder(homePath string, cf env.Finder, file stream.FileReader) Finder {
	return Finder{
		homePath:  homePath,
		ctxFinder: cf,
		file:      file,
	}
}

func (f Finder) Find(provider string) (Detail, error) {
	ctx, err := f.ctxFinder.Find()

	if err != nil {
		return Detail{}, err
	}
	if ctx.Current == "" {
		ctx.Current = env.DefaultEnv
	}

	cb, err := f.file.Read(File(f.homePath, ctx.Current, provider))
	if err != nil {
		errMsg := fmt.Sprintf(errNotFoundTemplate, provider)
		return Detail{Credential: Credential{}}, errors.New(prompt.Red(errMsg))
	}

	cred := &Detail{}
	if err := json.Unmarshal(cb, cred); err != nil {
		return Detail{}, err
	}
	return *cred, nil

}
