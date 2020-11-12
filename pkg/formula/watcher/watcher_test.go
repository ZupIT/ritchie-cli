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

package watcher

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/radovskyb/watcher"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/builder"
	"github.com/ZupIT/ritchie-cli/pkg/formula/repo"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/stream/streams"
)

func TestWatch(t *testing.T) {
	tmpDir := os.TempDir()
	workspacePath := filepath.Join(tmpDir, "ritchie-formulas-test-watcher")
	formulaPath := filepath.Join(tmpDir, "ritchie-formulas-test-watcher", "testing", "formula")
	ritHome := filepath.Join(os.TempDir(), ".my-rit-watcher")
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	treeGenerator := tree.NewGenerator(dirManager, fileManager)

	repoProviders := formula.NewRepoProviders()
	repoCreator := repo.NewCreator(ritHome, repoProviders, dirManager, fileManager)
	repoLister := repo.NewLister(ritHome, fileManager)
	repoWriter := repo.NewWriter(ritHome, fileManager)
	repoListWriteCreator := repo.NewListWriteCreator(repoLister, repoCreator, repoWriter)
	repoAdder := repo.NewAdder(ritHome, repoListWriteCreator, treeGenerator, fileManager)

	_ = dirManager.Remove(ritHome)
	_ = dirManager.Remove(workspacePath)
	_ = dirManager.Create(workspacePath)
	zipFile := filepath.Join("..", "..", "..", "testdata", "ritchie-formulas-test.zip")
	_ = streams.Unzip(zipFile, workspacePath)

	builderManager := builder.NewBuildLocal(ritHome, dirManager, repoAdder)
	sendMetric := func(commandExecutionTime float64, err ...string) {}

	watchManager := New(
		builderManager,
		dirManager,
		sendMetric,
	)

	go func() {
		watchManager.watcher.Wait()
		watchManager.watcher.TriggerEvent(watcher.Create, nil)
		watchManager.watcher.Error <- errors.New("error to watch formula")
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()

	workspace := formula.Workspace{
		Name: "repo-1",
		Dir:  workspacePath,
	}
	watchManager.Watch(formulaPath, workspace)

	hasRitchieHome := dirManager.Exists(ritHome)
	if !hasRitchieHome {
		t.Error("Watch build did not create the Ritchie home directory")
	}

	treeLocalFile := filepath.Join(ritHome, "repos", "local-repo-1", "tree.json")
	hasTreeLocalFile := fileManager.Exists(treeLocalFile)
	if !hasTreeLocalFile {
		t.Error("Watch build did not copy the tree local file")
	}

	formulaFiles := filepath.Join(ritHome, "repos", "local-repo-1", "testing", "formula", "bin")
	files, err := fileManager.List(formulaFiles)
	if err == nil && len(files) != 4 {
		t.Error("Watch build did not generate formulas files")
	}

	configFile := filepath.Join(ritHome, "repos", "local-repo-1", "testing", "formula", "config.json")
	hasConfigFile := fileManager.Exists(configFile)
	if !hasConfigFile {
		t.Error("Watch build did not copy formula config")
	}

}
