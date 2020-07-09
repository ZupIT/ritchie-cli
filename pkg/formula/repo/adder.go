package repo

import (
	"encoding/json"
	"fmt"
	"path"
	"sort"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/github"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/stream/streams"
)

const (
	reposDirName  = "repos"
	reposFileName = "repositories.json"
)

type AddManager struct {
	ritHome string
	github  github.Repositories
	tree    formula.TreeGenerator
	dir     stream.DirCreateListCopyRemover
	file    stream.FileWriteCreatorReadExistRemover
}

func NewAdder(
	ritHome string,
	github github.Repositories,
	tree formula.TreeGenerator,
	dir stream.DirCreateListCopyRemover,
	file stream.FileWriteCreatorReadExistRemover,
) AddManager {
	return AddManager{
		ritHome: ritHome,
		github:  github,
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
	repoInfo := github.NewRepoInfo(repo.Url, repo.Token)
	zipball, err := ad.github.Zipball(repoInfo, repo.Version)
	if err != nil {
		return err
	}

	defer zipball.Close()

	newRepoPath := path.Join(ad.ritHome, reposDirName, repo.Name)
	if err := ad.dir.Remove(newRepoPath); err != nil {
		return err
	}

	if err := ad.dir.Create(newRepoPath); err != nil {
		return err
	}

	zipFile := path.Join(newRepoPath, fmt.Sprintf("%s.zip", repo.Name))
	if err := ad.file.Create(zipFile, zipball); err != nil {
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
	exist := func() bool {
		for i := range repos {
			r := repos[i]
			if repo.Name == r.Name {
				repos[i].Priority = repo.Priority
				return true
			}
		}
		return false
	}

	if !exist() {
		repos = append(repos, repo)
	}

	for i := range repos {
		r := repos[i]
		if repo.Name != r.Name && r.Priority >= repo.Priority {
			repos[i].Priority++
		}
	}

	sort.Sort(repos)

	return repos
}
