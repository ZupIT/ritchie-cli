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

package tree

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/stream/streams"
)

func TestGenerate(t *testing.T) {
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	generator := NewGenerator(dirManager, fileManager)

	tmpDir := os.TempDir()
	workspacePath := fmt.Sprintf("%s/ritchie-formulas-test", tmpDir)
	ritHome := fmt.Sprintf("%s/.my-rit", os.TempDir())

	_ = dirManager.Remove(ritHome)
	_ = dirManager.Remove(workspacePath)
	_ = dirManager.Create(workspacePath)
	_ = streams.Unzip("../../../testdata/ritchie-formulas-test.zip", workspacePath)

	tree, err := generator.Generate(workspacePath)
	if err != nil {
		t.Error(err)
	}

	bytes, err := json.MarshalIndent(tree, "", "\t")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(bytes))
}
