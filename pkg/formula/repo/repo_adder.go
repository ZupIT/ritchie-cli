package repo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/http/headers"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/stream/streams"
)

const (
	repoDirPath  = "repo"
	repoFilePath = "repositories.json"
)

type AddManager struct {
	ritHome string
	client  *http.Client
	file    stream.FileWriteCreatorReadExistRemover
	dir     stream.DirCreateListCopyRemover
}

func NewAdder(ritHome string, client *http.Client, dir stream.DirCreateListCopyRemover, file stream.FileWriteCreatorReadExistRemover) AddManager {
	return AddManager{
		ritHome: ritHome,
		client:  client,
		dir:     dir,
		file:    file,
	}
}

func (ad AddManager) Add(repo formula.Repo) error {
	if err := ad.downloadRepo(repo); err != nil {
		return err
	}

	repos := formula.Repos{}
	repoPath := path.Join(ad.ritHome, repoDirPath, repoFilePath)
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

	return nil
}

func (ad AddManager) downloadRepo(repo formula.Repo) error {
	req, err := http.NewRequest(http.MethodGet, repo.ZipUrl, nil)
	if err != nil {
		return err
	}

	if repo.Token != "" {
		authToken := fmt.Sprintf("token %s", repo.Token)
		req.Header.Add(headers.Authorization, authToken)
	}

	req.Header.Add(headers.Accept, "application/vnd.github.v3+json")
	resp, err := ad.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	newRepoPath := path.Join(ad.ritHome, repoDirPath, repo.Name)
	if err := ad.dir.Remove(newRepoPath); err != nil {
		return err
	}

	if err := ad.dir.Create(newRepoPath); err != nil {
		return err
	}

	zipFile := path.Join(newRepoPath, fmt.Sprintf("%s.zip", repo.Name))
	if err := ad.file.Create(zipFile, resp.Body); err != nil {
		return err
	}

	if err := streams.Unzip(zipFile, newRepoPath); err != nil {
		return err
	}

	if err := ad.file.Remove(zipFile); err != nil {
		return err
	}

	dirs, err := ad.dir.List(newRepoPath, false)
	if err != nil {
		return err
	}

	src := path.Join(newRepoPath, dirs[0])
	if err := ad.dir.Copy(src, newRepoPath); err != nil {
		return err
	}

	if err := ad.dir.Remove(src); err != nil {
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
