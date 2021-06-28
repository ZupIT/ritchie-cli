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

	if tree.Version != Version {
		t.Fatalf("Generate(valid tree version) = %s; want %s", tree.Version, Version)
	}

	const cmdSize = 17
	if len(tree.Commands) != cmdSize {
		t.Fatalf("Generate(valid tree commands size) = %d; want %d", len(tree.Commands), cmdSize)
	}
}

func BenchmarkGenerator(b *testing.B) {
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

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = generator.Generate(workspacePath)
	}
}
