package repo

import (
	"encoding/json"
	"errors"
	"path"
	"sort"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const (
	repositoryDoNotExistError = "there is no repositories yet"
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

func (sm SetPriorityManager) SetPriority(repoName formula.RepoName, priority int) error {
	var repos formula.Repos
	repoPath := path.Join(sm.ritHome, reposDirName, reposFileName)
	if !sm.file.Exists(repoPath) {
		return errors.New(repositoryDoNotExistError)
	}
	read, err := sm.file.Read(repoPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(read, &repos); err != nil {
		return err
	}

	for i := range repos {
		if repoName == repos[i].Name {
			repos[i].Priority = priority
		} else if repos[i].Priority >= priority {
			repos[i].Priority++
		}
	}

	sort.Sort(repos)

	bytes, err := json.MarshalIndent(repos, "", "\t")
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
