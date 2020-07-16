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
