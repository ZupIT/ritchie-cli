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
	"errors"
	"os"
	"os/exec"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

var ErrBuildFormulaMakefile = errors.New("failed building formula with make, verify your repository")

type MakeManager struct {}

func NewBuildMake() formula.MakeBuilder {
	return MakeManager{}
}

func (ma MakeManager) Build(formulaPath string) error {
	if err := os.Chdir(formulaPath); err != nil {
		return err
	}
	cmd := exec.Command("make", "build")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return ErrBuildFormulaMakefile
	}

	return nil
}