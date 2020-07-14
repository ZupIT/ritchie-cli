package builder

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileextensions"
	"github.com/ZupIT/ritchie-cli/pkg/os/osutil"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const commonsDir = "commons"

var (
	msgBuildOnWindows = prompt.Yellow("This formula cannot be built on Windows.")
	ErrBuildOnWindows = errors.New(msgBuildOnWindows)
)

type Manager struct {
	ritHome string
	dir     stream.DirCreateListCopier
	file    stream.FileCopyExistLister
}

func New(ritHome string, dir stream.DirCreateListCopier, file stream.FileCopyExistLister) Manager {
	return Manager{ritHome: ritHome, dir: dir, file: file}
}

func (m Manager) Build(workspacePath, formulaPath string) error {
	formulaSrc := path.Join(formulaPath, "/src")
	if err := os.Chdir(formulaSrc); err != nil {
		return err
	}

	so := runtime.GOOS
	var cmd *exec.Cmd
	switch so {
	case osutil.Windows:
		winBuild := path.Join(formulaSrc, "build.bat")
		if !m.file.Exists(winBuild) {
			return ErrBuildOnWindows
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
			return errors.New(errMsg)
		}

		return err
	}

	formulaDestPath := m.formulaDestPath(formulaPath, workspacePath)

	if err := m.copyDist(formulaPath, formulaDestPath); err != nil {
		return err
	}

	if err := m.copyConfig(formulaPath, formulaDestPath); err != nil {
		return err
	}

	if err := m.copyTree(workspacePath); err != nil {
		return err
	}

	return nil
}

func (m Manager) copyDist(formulaPath, ritFormulaDistPath string) error {
	formulaDist := path.Join(formulaPath, "/dist") // /dist directory that contains built formula
	dirs, err := m.dir.List(formulaDist, false)
	if err != nil {
		return err
	}

	if err := m.dir.Create(ritFormulaDistPath); err != nil {
		return err
	}

	so := runtime.GOOS
	for _, dir := range dirs {
		if dir == so { // The formula is compiled and generates a binary by S.O, for example Golang, C...
			// Create formulaDistSODir for dist with the current S.O. example: "/dist/linux"
			formulaDistSODir := path.Join(formulaDist, "/", so)
			if err := m.dir.Copy(formulaDistSODir, ritFormulaDistPath); err != nil { // Copy formula dist by S.O. to ~/.rit/formulas/...
				return err
			}
			break
		}

		if dir == commonsDir { // The formula is interpreted and needs other files to run, for example, Java, Node...
			formulaCommonsDir := path.Join(formulaDist, "/", commonsDir)
			if err := m.dir.Copy(formulaCommonsDir, ritFormulaDistPath); err != nil { // Copy formula dist commons to ~/.rit/formulas/...
				return err
			}
			break
		}
	}

	return nil
}

func (m Manager) copyConfig(formulaPath string, distPath string) error {
	files, err := m.file.List(formulaPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.Contains(file, fileextensions.Json) {
			copyFile := path.Join(formulaPath, "/", file)
			distFile := path.Join(distPath, "/", file)
			if err := m.file.Copy(copyFile, distFile); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m Manager) copyTree(workspacePath string) error {
	copyFile := path.Join(workspacePath, "/tree/tree.json")
	destFile := path.Join(m.ritHome, "/repo/local/tree.json")
	destDir := path.Join(m.ritHome, "/repo/local")

	if err := m.dir.Create(destDir); err != nil {
		return err
	}

	if err := m.file.Copy(copyFile, destFile); err != nil {
		return err
	}

	return nil
}

func (m Manager) formulaDestPath(formulaPath, workspacePath string) string {
	dest := strings.ReplaceAll(formulaPath, workspacePath, "")
	return path.Join(m.ritHome, "/formulas/", dest)
}
