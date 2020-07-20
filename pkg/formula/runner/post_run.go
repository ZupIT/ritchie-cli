package runner

import (
	"fmt"
	"os/exec"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

type PostRunnerManager struct {
}

func NewPostRunner() PostRunnerManager {
	return PostRunnerManager{}
}

func (PostRunnerManager) PostRun(p formula.Setup, docker bool) error {
	if docker {
		if err := fileutil.RemoveFile(envFile); err != nil {
			return err
		}

		if err := removeContainer(p.ContainerId); err != nil {
			return err
		}
	}

	defer removeWorkDir(p.TmpDir)

	df, err := fileutil.ListNewFiles(p.BinPath, p.TmpDir)
	if err != nil {
		return err
	}

	if err = fileutil.MoveFiles(p.TmpDir, p.Pwd, df); err != nil {
		return err
	}

	return nil
}

func removeWorkDir(tmpDir string) {
	if err := fileutil.RemoveDir(tmpDir); err != nil {
		fmt.Sprintln("Error in remove dir")
	}
}

func removeContainer(imgName string) error {
	args := []string{"rm", imgName}
	cmd := exec.Command(dockerCmd, args...)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
