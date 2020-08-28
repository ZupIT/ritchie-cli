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

package runner

import (
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type PostRunnerManager struct {
	file stream.FileNewListMoveRemover
	dir  stream.DirRemover
}

func NewPostRunner(file stream.FileNewListMoveRemover, dir stream.DirRemover) PostRunnerManager {
	return PostRunnerManager{file: file, dir: dir}
}

func (po PostRunnerManager) PostRun(p formula.Setup, docker bool) error {
	if docker {
		if err := po.file.Remove(".env"); err != nil {
			return err
		}
	}

	defer po.removeWorkDir(p.TmpDir)

	df, err := po.file.ListNews(p.BinPath, p.TmpDir)
	if err != nil {
		return err
	}

	if err = po.file.Move(p.TmpDir, p.Pwd, df); err != nil {
		return err
	}

	return nil
}

func (po PostRunnerManager) removeWorkDir(tmpDir string) {
	if err := po.dir.Remove(tmpDir); err != nil {
		fmt.Sprintln("Error in remove dir")
	}
}
