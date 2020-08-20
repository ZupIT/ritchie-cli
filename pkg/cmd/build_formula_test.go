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

package cmd

import (
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula/builder"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/formula/watcher"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestNewBuildFormulaCmd(t *testing.T) {
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	treeGen := tree.NewGenerator(dirManager, fileManager)
	formulaLocalBuilder := builder.NewBuildLocal(os.TempDir(), dirManager, fileManager, treeGen)
	watchManager := watcher.New(formulaLocalBuilder, dirManager)

	cmd := NewBuildFormulaCmd(
		os.TempDir(),
		formulaLocalBuilder,
		workspaceForm{},
		watchManager,
		dirManager,
		inputTextMock{},
		inputListMock{},
		TutorialFinderMock{},
	)
	if cmd == nil {
		t.Errorf("NewBuildFormulaCmd got %v", cmd)
	}
}
