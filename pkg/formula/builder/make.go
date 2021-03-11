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
)

const msgMakeBuildErr = "failed building formula with Makefile, verify your repository"

var ErrBuildFormulaMakefile = errors.New(msgMakeBuildErr)

var _ formula.Builder = MakeManager{}

type MakeManager struct{}

func NewBuildMake() MakeManager {
	return MakeManager{}
}

func (ma MakeManager) Build(info formula.BuildInfo) error {
	pwd, _ := os.Getwd()
	defer os.Chdir(pwd) //nolint:errcheck

	if err := os.Chdir(info.FormulaPath); err != nil {
		return err
	}

	var stderr bytes.Buffer
	cmd := exec.Command("make", "build")
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf(errMsgFmt, ErrBuildFormulaMakefile, stderr.String())
	}

	if err := os.Chdir(pwd); err != nil {
		return err
	}

	return nil
}
