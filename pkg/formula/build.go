package formula

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
)

type BuilderManager struct {
	ritHome string
}

func NewBuilder(ritHome string) BuilderManager {
	return BuilderManager{ritHome: ritHome}
}

func (b BuilderManager) Build(workspacePath, formulaPath string) error {
	formulaSrc := fmt.Sprintf("%s/src", formulaPath)
	if err := os.Chdir(formulaSrc); err != nil { // cd formula src directory
		return err
	}

	cmd := exec.Command("make", "build") // make build formula
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	formulaDist := fmt.Sprintf("%s/dist", formulaPath) // cd formula dist directory, before build
	so := runtime.GOOS

	file, err := os.Open(formulaDist)
	if err != nil {
		return err
	}

	dirs, err := file.Readdirnames(0)
	if err != nil {
		return err
	}

	dist := strings.ReplaceAll(formulaPath, workspacePath, "")
	distPath := fmt.Sprintf("%s/formulas/%s", b.ritHome, dist)
	if err := fileutil.CreateDirIfNotExists(distPath, 0777); err != nil {
		return err
	}

	for _, dir := range dirs {
		if dir == so {
			path := fmt.Sprintf("%s/%s", formulaDist, so)
			if err := fileutil.CopyDirectory(path, distPath); err != nil {
				return err
			}
			break
		}

		if dir == "commons" {
			path := fmt.Sprintf("%s/%s", formulaDist, "commons")
			if err := fileutil.CopyDirectory(path, distPath); err != nil {
				return err
			}
			break
		}
	}

	if err := copyFormulaConfig(formulaPath, distPath); err != nil {
		return err
	}

	if err := b.copyTree(workspacePath); err != nil {
		return err
	}

	return nil
}

func copyFormulaConfig(formulaPath string, distPath string) error {
	file, err := os.Open(formulaPath)
	if err != nil {
		return err
	}

	files, err := file.Readdirnames(0)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.Contains(file, ".json") {
			copyFile := fmt.Sprintf("%s/%s", formulaPath, file)
			distFile := fmt.Sprintf("%s/%s", distPath, file)
			if err := fileutil.Copy(copyFile, distFile); err != nil {
				return err
			}
		}
	}

	return nil
}

func (b BuilderManager) copyTree(workspacePath string) error {
	copyFile := fmt.Sprintf("%s/tree/tree.json", workspacePath)
	distFile := fmt.Sprintf("%s/repo/local/tree.json",  b.ritHome)
	distDir := fmt.Sprintf("%s/repo/local",  b.ritHome)

	if err := fileutil.CreateDirIfNotExists(distDir, 0777); err != nil {
		return err
	}

	if err := fileutil.Copy(copyFile, distFile); err != nil {
		return err
	}

	return nil
}
