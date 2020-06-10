package formula

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const (
	srcPattern         = "%s/src"
	distPattern        = "%s/dist"
	formulaPathPattern = "%s/formulas/%s"
	defaultPattern     = "%s/%s"
	commonsDir         = "commons"
)

type BuilderManager struct {
	ritHome string
	dir     stream.DirCreateListCopier
	file    stream.FileListCopier
}

func NewBuilder(ritHome string, dir stream.DirCreateListCopier, file stream.FileListCopier) BuilderManager {
	return BuilderManager{ritHome: ritHome, dir: dir, file: file}
}

func (b BuilderManager) Build(workspacePath, formulaPath string) ([]byte, error) {
	formulaSrc := fmt.Sprintf(srcPattern, formulaPath)
	if err := os.Chdir(formulaSrc); err != nil { // cd formula src directory
		return nil, err
	}

	cmd := exec.Command("make", "build") // make build formula

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return stderr.Bytes(), err
	}

	dest := strings.ReplaceAll(formulaPath, workspacePath, "")
	ritFormulaDestPath := fmt.Sprintf(formulaPathPattern, b.ritHome, dest)

	if err := b.copyDist(formulaPath, ritFormulaDestPath); err != nil {
		return nil, err
	}

	if err := b.copyFormulaConfig(formulaPath, ritFormulaDestPath); err != nil {
		return nil, err
	}

	if err := b.copyTree(workspacePath); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b BuilderManager) copyDist(formulaPath, ritFormulaDistPath string) error {
	formulaDist := fmt.Sprintf(distPattern, formulaPath) // /dist directory that contains built formula
	dirs, err := b.dir.List(formulaDist, false)
	if err != nil {
		return err
	}

	if err := b.dir.Create(ritFormulaDistPath); err != nil {
		return err
	}

	so := runtime.GOOS
	for _, dir := range dirs {
		if dir == so { // The formula is compiled and generates a binary by S.O, for example Golang, C...
			// Create formulaDistSODir for dist with the current S.O. example: "/dist/linux"
			formulaDistSODir := fmt.Sprintf(defaultPattern, formulaDist, so)
			if err := b.dir.Copy(formulaDistSODir, ritFormulaDistPath); err != nil { // Copy formula dist by S.O. to ~/.rit/formulas/...
				return err
			}
			break
		}

		if dir == commonsDir { // The formula is interpreted and needs other files to run, for example, Java, Node...
			formulaCommonsDir := fmt.Sprintf(defaultPattern, formulaDist, commonsDir)
			if err := b.dir.Copy(formulaCommonsDir, ritFormulaDistPath); err != nil { // Copy formula dist commons to ~/.rit/formulas/...
				return err
			}
			break
		}
	}

	return nil
}

func (b BuilderManager) copyFormulaConfig(formulaPath string, distPath string) error {
	files, err := b.file.List(formulaPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.Contains(file, ".json") {
			copyFile := fmt.Sprintf(defaultPattern, formulaPath, file)
			distFile := fmt.Sprintf(defaultPattern, distPath, file)
			if err := b.file.Copy(copyFile, distFile); err != nil {
				return err
			}
		}
	}

	return nil
}

func (b BuilderManager) copyTree(workspacePath string) error {
	copyFile := fmt.Sprintf("%s/tree/tree.json", workspacePath)
	destFile := fmt.Sprintf("%s/repo/local/tree.json", b.ritHome)
	destDir := fmt.Sprintf("%s/repo/local", b.ritHome)

	if err := b.dir.Create(destDir); err != nil {
		return err
	}

	if err := b.file.Copy(copyFile, destFile); err != nil {
		return err
	}

	return nil
}
