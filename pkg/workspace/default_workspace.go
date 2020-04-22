package workspace

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"

	"github.com/gofrs/flock"
)

const (
	repoDir  = "/repo"
	repoFile = "/repositories.json"
)

type DefaultChecker struct {
	ritchieHome string
	dir         stream.DirCreater
	file        stream.FileWriteReadExister
}

func NewChecker(ritchieHome string, dir stream.DirCreater, file stream.FileWriteReadExister) DefaultChecker {
	return DefaultChecker{
		ritchieHome: ritchieHome,
		dir:         dir,
		file:        file,
	}
}

func (d DefaultChecker) Check() error {
	dirRepo := fmt.Sprintf("%s%s", d.ritchieHome, repoDir)
	repoFile := fmt.Sprintf("%s%s", dirRepo, repoFile)

	if err := d.dir.Create(d.ritchieHome); err != nil {
		return err
	}

	if err := d.dir.Create(dirRepo); err != nil {
		return err
	}

	if d.file.Exists(repoFile) {
		return nil
	}

	lock := flock.New(strings.Replace(repoFile, filepath.Ext(repoFile), ".lock", 1))
	lockCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	locked, err := lock.TryLockContext(lockCtx, time.Second)
	if locked {
		defer lock.Unlock()
	}

	if err != nil {
		return err
	}

	b, err := json.Marshal(formula.RepositoryFile{})
	if err != nil {
		return err
	}
	d.file.Write(repoFile, b)

	return nil
}
