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

package rcontext

import (
	"encoding/json"
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type FindManager struct {
	CtxFile string
	File    stream.FileReadExister
}

func NewFinder(homePath string, file stream.FileReadExister) FindManager {
	return FindManager{
		CtxFile: fmt.Sprintf(ContextPath, homePath),
		File:    file,
	}
}

func (f FindManager) Find() (ContextHolder, error) {
	ctxHolder := ContextHolder{}

	if !f.File.Exists(f.CtxFile) {
		return ctxHolder, nil
	}

	file, err := f.File.Read(f.CtxFile)
	if err != nil {
		return ctxHolder, err
	}

	if err := json.Unmarshal(file, &ctxHolder); err != nil {
		return ctxHolder, err
	}

	return ctxHolder, nil
}
