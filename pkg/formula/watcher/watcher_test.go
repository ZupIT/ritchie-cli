package watcher

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/radovskyb/watcher"

	"github.com/ZupIT/ritchie-cli/pkg/formula/builder"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/stream/streams"
)

func TestWatch(t *testing.T) {
	tmpDir := os.TempDir()
	workspacePath := fmt.Sprintf("%s/ritchie-formulas-test-watcher", tmpDir)
	formulaPath := fmt.Sprintf("%s/ritchie-formulas-test-watcher/testing/formula", tmpDir)
	ritHome := fmt.Sprintf("%s/.my-rit-watcher", os.TempDir())
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	treeGenerator := tree.NewGenerator(dirManager, fileManager)

	_ = dirManager.Remove(ritHome)
	_ = dirManager.Remove(workspacePath)
	_ = dirManager.Create(workspacePath)
	_ = streams.Unzip("../../../testdata/ritchie-formulas-test.zip", workspacePath)

	builderManager := builder.New(ritHome, dirManager, fileManager, treeGenerator)

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

	treeLocalFile := fmt.Sprintf("%s/repos/local/tree.json", ritHome)
	hasTreeLocalFile := fileManager.Exists(treeLocalFile)
	if !hasTreeLocalFile {
		t.Error("Watch build did not copy the tree local file")
	}

	formulaFiles := fmt.Sprintf("%s/repos/local/testing/formula/bin", ritHome)
	files, err := fileManager.List(formulaFiles)
	if err == nil && len(files) != 2 {
		t.Error("Watch build did not generate formulas files")
	}

	configFile := fmt.Sprintf("%s/repos/local/testing/formula/config.json", ritHome)
	hasConfigFile := fileManager.Exists(configFile)
	if !hasConfigFile {
		t.Error("Watch build did not copy formula config")
	}
}
