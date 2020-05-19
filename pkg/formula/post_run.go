package formula

import (
	"fmt"
	"os/exec"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
)

func PostRun(p Setup, docker bool) error {
	if docker {
		if err := fileutil.RemoveFile(envFile); err != nil {
			return err
		}

		if err := removeContainer(p.containerId); err != nil {
			return err
		}
	}

	defer removeWorkDir(p.tmpDir)

	df, err := fileutil.ListNewFiles(p.binPath, p.tmpBinDir)
	if err != nil {
		return err
	}

	if err = fileutil.MoveFiles(p.tmpBinDir, p.pwd, df); err != nil {
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