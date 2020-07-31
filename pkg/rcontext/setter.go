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
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/slice/sliceutil"
)

type SetterManager struct {
	ctxFile string
	finder  Finder
}

func NewSetter(homePath string, f Finder) Setter {
	return SetterManager{ctxFile: fmt.Sprintf(ContextPath, homePath), finder: f}
}

func (s SetterManager) Set(ctx string) (ContextHolder, error) {
	ctxHolder, err := s.finder.Find()
	if err != nil {
		return ContextHolder{}, err
	}

	ctxHolder.Current = strings.ReplaceAll(ctx, DefaultCtx, "")
	if ctx != DefaultCtx {
		if ctxHolder.All == nil {
			ctxHolder.All = make([]string, 0)
		}

		if !sliceutil.Contains(ctxHolder.All, ctx) {
			ctxHolder.All = append(ctxHolder.All, ctx)
		}
	}

	b, err := json.Marshal(&ctxHolder)
	if err != nil {
		return ContextHolder{}, err
	}
	if err := fileutil.WriteFilePerm(s.ctxFile, b, 0600); err != nil {
		return ContextHolder{}, err
	}

	return ctxHolder, nil
}
