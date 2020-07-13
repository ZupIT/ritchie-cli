package repo

import (
	"encoding/json"
	"path"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type DeleteManager struct {
	ritHome string
	file    stream.FileWriteReadExister
	dir     stream.DirRemover
}

func NewDeleter(ritHome string, file stream.FileWriteReadExister, dir stream.DirRemover) DeleteManager {
	return DeleteManager{
		ritHome: ritHome,
		file:    file,
		dir:     dir,
	}
}

func (dm DeleteManager) Delete(repoName string) error {
	if err := dm.deleteRepoDir(repoName); err != nil {
		return err
	}
	if err := dm.deleteFromReposFile(repoName); err != nil {
		return err
	}
	return nil
}

func (dm DeleteManager) deleteRepoDir(repoName string) error {
	path := path.Join(dm.ritHome, reposDirName, repoName)
	if err := dm.dir.Remove(path); err != nil {
		return err
	}
	return nil
}

func (dm DeleteManager) deleteFromReposFile(repoName string) error {
	repos := formula.Repos{}

	path := path.Join(dm.ritHome, reposDirName, reposFileName)
	file, err := dm.file.Read(path)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(file, &repos); err != nil {
		return err
	}

	var idx int
	for i := range repos {
		if repos[i].Name == repoName {
			idx = i
			break
		}
	}
	repos = append(repos[:idx], repos[idx+1:]...)

	newFile, err := json.MarshalIndent(repos, "", "\t")
	if err != nil {
		return err
	}

	if err = dm.file.Write(path, newFile); err != nil {
		return err
	}

	return nil
}
