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

package prompt

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const (
	docsDir                            = "docs"
	srcDir                             = "src"
	questionSelectFormulaGroup         = "Select a formula or group: "
	optionOtherFormula                 = "Another formula"
	InputFormulaQuestionFoundedFormula = "we found a formula, do you can select that: "
)

type InputFormula struct {
	inList    InputList
	directory stream.DirListChecker
}

func NewInputFormula(inList InputList, directory stream.DirListChecker) InputFormula {
	return InputFormula{inList: inList, directory: directory}
}

func (i *InputFormula) Select(dir string, currentFormula string) (string, error) {
	formula := ""

	group, err := i.readFormulas(dir, currentFormula)
	if len(group) > 0 {
		formula = strings.Join(group, " ")
	}

	return formula, err
}

func (i *InputFormula) readFormulas(dir string, currentFormula string) ([]string, error) {
	dirs, err := i.directory.List(dir, false)
	if err != nil {
		return nil, err
	}

	dirs = removeFromArray(dirs, docsDir)

	var groups []string
	var formulaOptions []string
	var response string

	if isFormula(dirs) {
		if !hasFormulaInDir(dirs) {
			return groups, nil
		}

		formulaOptions = append(formulaOptions, currentFormula, optionOtherFormula)

		response, err = i.inList.List(InputFormulaQuestionFoundedFormula, formulaOptions)
		if err != nil {
			return nil, err
		}
		if response == currentFormula {
			return groups, nil
		}
		dirs = removeFromArray(dirs, srcDir)
	}

	selected, err := i.inList.List(questionSelectFormulaGroup, dirs)
	if err != nil {
		return nil, err
	}

	newFormulaSelected := fmt.Sprintf("%s %s", currentFormula, selected)

	var aux []string
	aux, err = i.readFormulas(filepath.Join(dir, selected), newFormulaSelected)
	if err != nil {
		return nil, err
	}

	aux = append([]string{selected}, aux...)
	groups = append(groups, aux...)

	return groups, nil
}

func removeFromArray(ss []string, r string) []string {
	for i, s := range ss {
		if s == r {
			return append(ss[:i], ss[i+1:]...)
		}
	}
	return ss
}

func isFormula(dirs []string) bool {
	for _, dir := range dirs {
		if dir == srcDir {
			return true
		}
	}

	return false
}

func hasFormulaInDir(dirs []string) bool {
	dirs = removeFromArray(dirs, docsDir)
	dirs = removeFromArray(dirs, srcDir)

	return len(dirs) > 0
}
