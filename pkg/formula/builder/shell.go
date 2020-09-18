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
	"path/filepath"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

const (
	buildSh          = "build.sh"
	msgShellBuildErr = "failed building formula with shell script, verify your repository"
)

var ErrBuildFormulaShell = errors.New(msgShellBuildErr)

var _ formula.ShellBuilder = ShellManager{}

type ShellManager struct {
}

func NewBuildShell() formula.ShellBuilder {
	return ShellManager{}
}

func (sh ShellManager) Build(formulaPath string) error {
	if err := os.Chdir(formulaPath); err != nil {
		return err
	}
	var stderr bytes.Buffer
	execFile := filepath.Join(formulaPath, buildSh)
	cmd := exec.Command(execFile)
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf(errMsgFmt, ErrBuildFormulaShell, stderr.String())
	}

	return nil
}
