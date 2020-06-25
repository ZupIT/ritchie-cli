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

	df, err := fileutil.ListNewFiles(p.BinPath, p.TmpBinDir)
	if err != nil {
		return err
	}

	if err = fileutil.MoveFiles(p.TmpBinDir, p.Pwd, df); err != nil {
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
	args := []string{dockerRemoveCmd, imgName}
	cmd := exec.Command(docker, args...)

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}
