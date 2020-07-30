package runner

import (
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type PostRunnerManager struct {
	file stream.FileMoveRemover
	dir  stream.DirRemover
}

func NewPostRunner(file stream.FileMoveRemover, dir stream.DirRemover) PostRunnerManager {
	return PostRunnerManager{file: file, dir: dir}
}

func (po PostRunnerManager) PostRun(p formula.Setup, docker bool) error {
	if docker {
		if err := po.file.Remove(envFile); err != nil {
			return err
		}
	}

	defer po.removeWorkDir(p.TmpDir)

	df, err := fileutil.ListNewFiles(p.BinPath, p.TmpDir)
	if err != nil {
		return err
	}

	if err = po.file.Move(p.TmpDir, p.Pwd, df); err != nil {
		return err
	}

	return nil
}

func (po PostRunnerManager) removeWorkDir(tmpDir string) {
	if err := po.dir.Remove(tmpDir); err != nil {
		fmt.Sprintln("Error in remove dir")
	}
}