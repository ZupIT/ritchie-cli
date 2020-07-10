package repo

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type UpdateManager struct {
	ritHome string
	repo    formula.RepositoryListCreator
	tree    formula.TreeGenerator
	file    stream.FileWriter
}

func NewUpdater(ritHome string,
	repo formula.RepositoryListCreator,
	tree formula.TreeGenerator,
	file stream.FileWriter,
) UpdateManager {
	return UpdateManager{
		ritHome: ritHome,
		repo:    repo,
		tree:    tree,
		file:    file,
	}
}

func (up UpdateManager) Update(name formula.RepoName, version formula.RepoVersion) error {
	repos, err := up.repo.List()
	if err != nil {
		return err
	}

	var repo *formula.Repo
	for i := range repos {
		if name == repos[i].Name {
			repo = &repos[i]
			break
		}
	}

	if repo == nil {
		return fmt.Errorf("repository name %q was not found", name)
	}

	repo.Version = version

	if err := up.repo.Create(*repo); err != nil {
		return err
	}

	repoFilePath := path.Join(up.ritHome, reposDirName, reposFileName)
	if err := up.saveRepo(repoFilePath, repos); err != nil {
		return err
	}

	repoPath := path.Join(up.ritHome, reposDirName, name.String())
	tree, err := up.tree.Generate(repoPath)
	if err != nil {
		return err
	}

	treeFilePath := path.Join(repoPath, "tree.json")
	bytes, err := json.MarshalIndent(tree, "", "\t")
	if err != nil {
		return err
	}

	if err := up.file.Write(treeFilePath, bytes); err != nil {
		return err
	}

	return nil
}

func (up UpdateManager) saveRepo(repoPath string, repos formula.Repos) error {
	bytes, err := json.MarshalIndent(repos, "", "\t")
	if err != nil {
		return err
	}

	if err := up.file.Write(repoPath, bytes); err != nil {
		return err
	}

	return nil
}
