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
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"strings"
)

type RemoveManager struct {
	ctxFile string
	finder  Finder
}

func NewRemover(homePath string, f Finder) RemoveManager {
	return RemoveManager{ctxFile: fmt.Sprintf(ContextPath, homePath), finder: f}
}

func (r RemoveManager) Remove(ctx string) (ContextHolder, error) {
	ctxHolder, err := r.finder.Find()
	if err != nil {
		return ContextHolder{}, err
	}

	ctx = strings.ReplaceAll(ctx, CurrentCtx, "")
	if ctxHolder.Current == ctx {
		ctxHolder.Current = ""
	}

	for i, context := range ctxHolder.All {
		if ctx == context {
			ctxHolder.All = append(ctxHolder.All[:i], ctxHolder.All[i+1:]...)
			break
		}
	}

	b, err := json.Marshal(&ctxHolder)
	if err != nil {
		return ContextHolder{}, err
	}
	if err := fileutil.WriteFilePerm(r.ctxFile, b, 0600); err != nil {
		return ContextHolder{}, err
	}

	return ctxHolder, nil
}
