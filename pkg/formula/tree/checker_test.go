package tree

import (
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestChecker(t *testing.T) {
	t.Run("Running checker test", func(t *testing.T) {
		treeChecker := NewChecker(stream.DirManager{}, stream.FileManager{})
		treeChecker.CheckCommands()
	})
}
