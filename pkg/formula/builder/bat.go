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
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const (
	buildBat       = "build.bat"
	msgBatBuildErr = "Check if you have all the requirements to execute this formula."
	errMsgFmt      = "\nError building formula: "
)

var (
	ErrBuildFormulaBuildBat = errors.New(msgBatBuildErr)
	msgBuildOnWindows       = prompt.Yellow("This formula cannot be built on Windows.")
	ErrBuildOnWindows       = errors.New(msgBuildOnWindows)
)

var _ formula.Builder = BatManager{}

type BatManager struct {
	file stream.FileExister
}

func NewBuildBat(file stream.FileExister) BatManager {
	return BatManager{file: file}
}

func (ba BatManager) Build(info formula.BuildInfo) error {
	if err := os.Chdir(info.FormulaPath); err != nil {
		return err
	}

	if !ba.file.Exists(buildBat) {
		return ErrBuildOnWindows
	}

	var stderr bytes.Buffer
	cmd := exec.Command(buildBat)
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return errors.New(
			fmt.Sprint(
				prompt.Red(ErrBuildFormulaBuildBat.Error()),
				errMsgFmt+prompt.Red(stderr.String()),
			))
	}

	return nil
}
