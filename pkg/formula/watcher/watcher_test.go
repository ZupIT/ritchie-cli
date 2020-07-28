package watcher

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/radovskyb/watcher"

	"github.com/ZupIT/ritchie-cli/pkg/formula/builder"
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

	_ = dirManager.Remove(ritHome)
	_ = dirManager.Remove(workspacePath)
	_ = dirManager.Create(workspacePath)
	zipFile := filepath.Join("..", "..", "..", "testdata", "ritchie-formulas-test.zip")
	_ = streams.Unzip(zipFile, workspacePath)

	builderManager := builder.NewBuildLocal(ritHome, dirManager, fileManager, treeGenerator)

	watchManager := New(builderManager, dirManager)

	go func() {
		watchManager.watcher.Wait()
		watchManager.watcher.TriggerEvent(watcher.Create, nil)
		watchManager.watcher.Error <- errors.New("error to watch formula")
		watchManager.watcher.Close()
	}()

	watchManager.Watch(workspacePath, formulaPath)

	hasRitchieHome := dirManager.Exists(ritHome)
	if !hasRitchieHome {
		t.Error("Watch build did not create the Ritchie home directory")
	}

	treeLocalFile := filepath.Join(ritHome, "repos", "local", "tree.json")
	hasTreeLocalFile := fileManager.Exists(treeLocalFile)
	if !hasTreeLocalFile {
		t.Error("Watch build did not copy the tree local file")
	}

	formulaFiles := filepath.Join(ritHome, "repos", "local", "testing", "formula", "bin")
	files, err := fileManager.List(formulaFiles)
	if err == nil && len(files) != 4 {
		t.Error("Watch build did not generate formulas files")
	}

	configFile := filepath.Join(ritHome, "repos", "local", "testing", "formula", "config.json")
	hasConfigFile := fileManager.Exists(configFile)
	if !hasConfigFile {
		t.Error("Watch build did not copy formula config")
	}
}
