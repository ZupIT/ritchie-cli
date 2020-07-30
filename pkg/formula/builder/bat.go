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

const buildFile = "build.bat"

var ErrBuildFormulaBuildBat = errors.New("failed building formula with build.bat, verify your repository")

type BatManager struct{}

func NewBuildBat() formula.BatBuilder {
	return BatManager{}
}

func (ba BatManager) Build(formulaPath string) error {
	if err := os.Chdir(formulaPath); err != nil {
		return err
	}

	cmd := exec.Command(buildFile)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return ErrBuildFormulaBuildBat
	}

	return nil
}
