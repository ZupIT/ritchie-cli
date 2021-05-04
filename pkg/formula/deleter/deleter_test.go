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

package deleter

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/stream/streams"
)

func TestDeleterCmd(t *testing.T) {
	home := filepath.Join(os.TempDir(), "rit_test-Deleter")
	ritHome := filepath.Join(home, ".rit")

	defer os.RemoveAll(home)

	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	treeGen := tree.NewGenerator(dirManager, fileManager)

	reposPath := filepath.Join(ritHome, "repos")
	repoPathLocalDefault := filepath.Join(reposPath, "local-default")
	repoPathWS := filepath.Join(home, "ritchie-formulas-local")

	tests := []struct {
		name string
		in   formula.Delete
		want error
	}{
		{
			name: "success on prompt input",
			in: formula.Delete{
				GroupsFormula: []string{"testing", "formula"},
				Workspace: formula.Workspace{
					Name: formula.DefaultWorkspaceName,
					Dir:  filepath.Join(home, formula.DefaultWorkspaceDir),
				},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.RemoveAll(home)
			_ = dirManager.Remove(ritHome)
			createSaved := func(path string) {
				_ = dirManager.Remove(path)
				_ = dirManager.Create(path)
			}
			createSaved(repoPathLocalDefault)
			createSaved(repoPathWS)

			zipFile := filepath.Join("..", "..", "testdata", "ritchie-formulas-test.zip")
			zipRepositories := filepath.Join("..", "..", "testdata", "repositories.zip")
			zipTree := filepath.Join("..", "..", "testdata", "tree.zip")
			_ = streams.Unzip(zipRepositories, reposPath)
			_ = streams.Unzip(zipFile, repoPathLocalDefault)
			_ = streams.Unzip(zipTree, repoPathLocalDefault)
			_ = streams.Unzip(zipFile, repoPathWS)

			createTree(ritHome, repoPathWS, treeGen, fileManager)

			deleter := NewDeleter(dirManager, fileManager, treeGen, ritHome)

			got := deleter.Delete(tt.in)

			assert.Equal(t, tt.want, got)
		})
	}

}

func createTree(ritHome, ws string, tg formula.TreeGenerator, fm stream.FileWriteRemover) {
	localTree, _ := tg.Generate(ws)

	jsonString, _ := json.MarshalIndent(localTree, "", "\t")
	pathLocalTreeJSON := filepath.Join(ritHome, "repos", "local-default", "tree.json")
	_ = ioutil.WriteFile(pathLocalTreeJSON, jsonString, os.ModePerm)
}
