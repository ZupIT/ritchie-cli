package repo

import (
	"encoding/json"
	"net/http"
	"path"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const repositoriesPath = "/repo/repositories.json"

type AddManager struct {
	ritHome string
	client  *http.Client
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

	if repo.Current { // If the new repo has been set as current, you must set other repos to be non-current
		repos = unsetCurrentRepo(repos)
	}

	repos[repo.Name] = repo
	if err := ad.saveRepo(repoPath, repos); err != nil {
		return err
	}

	ad.downloadRepo()

	return nil
}

func (ad AddManager) downloadRepo(repo formula.Repo) error {
	req, err := http.NewRequest(http.MethodGet, repo.ZipUrl)
	if err != nil {
		return err
	}



	return nil
}

func (ad AddManager) saveRepo(repoPath string, repos formula.Repos) error {
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

func unsetCurrentRepo(repos formula.Repos) formula.Repos {
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
