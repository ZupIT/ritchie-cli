package repo

import (
	"encoding/json"
	"path"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const repositoriesPath = "/repo/repositories.json"

type AddManager struct {
	ritHome string
	file    stream.FileWriteReadExister
	dir     stream.DirCreater
}

func NewAdder(ritHome string, dir stream.DirCreater, file stream.FileWriteReadExister) AddManager {
	return AddManager{
		ritHome: ritHome,
		dir:     dir,
		file:    file,
	}
}

func (ad AddManager) Add(repo formula.Repo) error {
	repos := formula.Repos{}
	repoPath := path.Join(ad.ritHome, repositoriesPath)
	if ad.file.Exists(repoPath) {
		read, err := ad.file.Read(repoPath)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(read, &repos); err != nil {
			return err
		}
	}

	if repo.Current {
		repos = unsetCurrent(repos)
	}

	repos[repo.Name] = repo
	bytes, err := json.Marshal(repos)
	if err != nil {
		return err
	}

	dirPath := path.Dir(repoPath)
	if err := ad.dir.Create(dirPath); err != nil {
		return err
	}

	if err := ad.file.Write(repoPath, bytes); err != nil {
		return err
	}

	return nil
}

func unsetCurrent(repos formula.Repos) formula.Repos {
	for k := range repos {
		repo := repos[k]
		if repo.Current {
			repo.Current = false
			repos[k] = repo
			break
		}
	}

	return repos
}
