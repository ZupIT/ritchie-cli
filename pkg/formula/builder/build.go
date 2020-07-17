package builder

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/os/osutil"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const localRepoDir = "/repos/local"

var (
	msgBuildOnWindows = prompt.Yellow("This formula cannot be built on Windows.")
	ErrBuildOnWindows = errors.New(msgBuildOnWindows)
)

type Manager struct {
	ritHome string
	dir     stream.DirCreateListCopier
	file    stream.FileCopyExistListerWriter
	tree    formula.TreeGenerator
}

func New(ritHome string, dir stream.DirCreateListCopier, file stream.FileCopyExistListerWriter, tree formula.TreeGenerator) Manager {
	return Manager{ritHome: ritHome, dir: dir, file: file, tree: tree}
}

func (m Manager) Build(workspacePath, formulaPath string) error {

	dest := path.Join(m.ritHome, localRepoDir)

	if err := m.dir.Create(dest); err != nil {
		return err
	}

	if err := m.copyWorkSpace(workspacePath, dest); err != nil {
		return err
	}

	if err := m.generateTree(dest); err != nil {
		return err
	}

	err, done := m.buildFormulaBin(workspacePath, formulaPath, dest)
	if done {
		return err
	}

	return nil
}

func (m Manager) buildFormulaBin(workspacePath, formulaPath, dest string) (error, bool) {
	formulaSrc := strings.ReplaceAll(formulaPath, workspacePath, dest)
	if err := os.Chdir(formulaSrc); err != nil {
		return err, true
	}

	so := runtime.GOOS
	var cmd *exec.Cmd
	switch so {
	case osutil.Windows:
		winBuild := path.Join(formulaSrc, "build.bat")
		if !m.file.Exists(winBuild) {
			return ErrBuildOnWindows, true
		}
		cmd = exec.Command(winBuild)
	default:
		cmd = exec.Command("make", "build")
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if stderr.Bytes() != nil {
			errMsg := fmt.Sprintf("Build error: \n%s \n%s", stderr.String(), err)
			return errors.New(errMsg), true
		}

		return err, true
	}
	return nil, false
}

func (m Manager) generateTree(dest string) error {
	tree, err := m.tree.Generate(dest)
	if err != nil {
		return err
	}

	treeFilePath := path.Join(dest, "tree.json")
	treeIndented, err := json.MarshalIndent(tree, "", "\t")
	if err != nil {
		return err
	}

	if err := m.file.Write(treeFilePath, treeIndented); err != nil {
		return err
	}
	return nil
}

func (m Manager) copyWorkSpace(workspacePath string, dest string) error {
	return m.dir.Copy(workspacePath, dest)
}
