package repo

import (
	"encoding/json"
	"path"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type SetPriorityManager struct {
	ritHome string
	file    stream.FileWriteReadExister
	dir     stream.DirCreater
}

func NewPrioritySetter(ritHome string, file stream.FileWriteReadExister, dir stream.DirCreater) SetPriorityManager {
	return SetPriorityManager{
		ritHome: ritHome,
		file:    file,
		dir:     dir,
	}
}

func (sm SetPriorityManager) SetPriority(repo formula.Repo, priority int) error {
	var repos formula.RepoFile
	repoPath := path.Join(sm.ritHome, repositoriesPath)
	if sm.file.Exists(repoPath) {
		read, err := sm.file.Read(repoPath)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(read, &repos); err != nil {
			return err
		}
	}

	for idx, r := range repos.Values {
		if r.Name == repo.Name {
			repos.Values[idx].Priority = priority
			break
		}
	}

	bytes, err := json.Marshal(repos)
	if err != nil {
		return err
	}

	dirPath := path.Dir(repoPath)
	if err := sm.dir.Create(dirPath); err != nil {
		return err
	}

	if err := sm.file.Write(repoPath, bytes); err != nil {
		return err
	}

	return nil
}