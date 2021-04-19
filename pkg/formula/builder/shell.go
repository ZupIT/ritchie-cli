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
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

const (
	buildSh          = "build.sh"
	msgShellBuildErr = "Check if you have all the requirements to execute this formula."
)

var ErrBuildFormulaShell = errors.New(msgShellBuildErr)

var _ formula.Builder = ShellManager{}

type ShellManager struct {
}

func NewBuildShell() ShellManager {
	return ShellManager{}
}

func (sh ShellManager) Build(info formula.BuildInfo) error {
	if err := os.Chdir(info.FormulaPath); err != nil {
		return err
	}
	var stderr bytes.Buffer
	execFile := filepath.Join(info.FormulaPath, buildSh)
	cmd := exec.Command(execFile) //nolint:gosec
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return errors.New(
			fmt.Sprint(
				prompt.Red(ErrBuildFormulaShell.Error()),
				errMsgFmt+ prompt.Red(stderr.String()),
			))
	}

	return nil
}
