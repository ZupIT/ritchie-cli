package workspace

import (
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
)

type DefaultChecker struct {
	ritchieHome string
}

func NewChecker(ritchieHome string) DefaultChecker {
	return DefaultChecker{ritchieHome: ritchieHome}
}

func (d DefaultChecker) Check() error {
	/*dirRepo := fmt.Sprintf("%s%s", d.ritchieHome, repoDir)
	repoFile := fmt.Sprintf("%s%s", dirRepo, repoFile)*/

	if err := fileutil.CreateDirIfNotExists(d.ritchieHome, 0755); err != nil {
		return err
	}

	/*if err := fileutil.CreateDirIfNotExists(dirRepo, 0755); err != nil {
		return err
	}

	if fileutil.Exists(repoFile) {
		return nil
	}

	lock := flock.New(strings.Replace(repoFile, filepath.Ext(repoFile), ".lock", 1))
	lockCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	locked, err := lock.TryLockContext(lockCtx, time.Second)
	if locked {
		defer func() {
			err := lock.Unlock()
			if err != nil {
				fmt.Sprintln("Error in Unlock")
				return
			}
		}()
	}

	if err != nil {
		return err
	}

	b, err := json.Marshal(formula.RepositoryFile{})
	if err != nil {
		return err
	}
	err = fileutil.WriteFile(repoFile, b)
	if err != nil {
		return err
	}*/

	return nil
}
