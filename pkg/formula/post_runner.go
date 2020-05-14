package formula

import (
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
)

func (d DefaultRunner) PostRun(p RunData) error {
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
