package repo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"sort"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/http/headers"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/stream/streams"
)

const (
	reposDirName  = "repos"
	reposFileName = "repositories.json"
)

type AddManager struct {
	ritHome string
	client  *http.Client
	tree    formula.TreeGenerator
	dir     stream.DirCreateListCopyRemover
	file    stream.FileWriteCreatorReadExistRemover
}

func NewAdder(
	ritHome string,
	client *http.Client,
	tree formula.TreeGenerator,
	dir stream.DirCreateListCopyRemover,
	file stream.FileWriteCreatorReadExistRemover,
) AddManager {
	return AddManager{
		ritHome: ritHome,
		client:  client,
		tree:    tree,
		dir:     dir,
		file:    file,
	}
}

func (ad AddManager) Add(repo formula.Repo) error {
	if err := ad.downloadRepo(repo); err != nil {
		return err
	}

	repos := formula.Repos{}
	repoPath := path.Join(ad.ritHome, reposDirName, reposFileName)
	if ad.file.Exists(repoPath) {
		read, err := ad.file.Read(repoPath)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(read, &repos); err != nil {
			return err
		}
	}

	repos = setPriority(repo, repos)

	if err := ad.saveRepo(repoPath, repos); err != nil {
		return err
	}

	newRepoPath := path.Join(ad.ritHome, reposDirName, repo.Name)

	tree, err := ad.tree.Generate(newRepoPath)
	if err != nil {
		return err
	}

	treeFilePath := path.Join(newRepoPath, "tree.json")
	bytes, err := json.MarshalIndent(tree, "", "\t")
	if err != nil {
		return err
	}

	if err := ad.file.Write(treeFilePath, bytes); err != nil {
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

	newRepoPath := path.Join(ad.ritHome, reposDirName, repo.Name)
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
	bytes, err := json.MarshalIndent(repos, "", "\t")
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

func setPriority(repo formula.Repo, repos formula.Repos) formula.Repos {
	repos = append(repos, repo)

	for i := range repos {
		r := repos[i]
		if repo.Name != r.Name && r.Priority >= repo.Priority {
			repos[i].Priority++
		}
	}

	sort.Sort(repos)

	return repos
}
