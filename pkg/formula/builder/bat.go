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

package builder

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const (
	buildFile      = "build.bat"
	msgBatBuildErr = "failed building formula with build.bat, verify your repository"
	errMsgFmt      = `%s
More about error: %s`
)

var ErrBuildFormulaBuildBat = errors.New(msgBatBuildErr)

type BatManager struct {
	file stream.FileExister
}

func NewBuildBat(file stream.FileExister) formula.BatBuilder {
	return BatManager{file: file}
}

func (ba BatManager) Build(formulaPath string) error {
	if err := os.Chdir(formulaPath); err != nil {
		return err
	}

	if !ba.file.Exists(buildFile) {
		return ErrBuildOnWindows
	}

	var stderr bytes.Buffer
	cmd := exec.Command(buildFile)
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf(errMsgFmt, ErrBuildFormulaBuildBat, stderr.String())
	}

	return nil
}
