package repo

import (
	"encoding/json"
	"path"
	"sort"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type ListManager struct {
	ritHome string
	file    stream.FileReadExister
}

func NewLister(ritHome string, file stream.FileReadExister) ListManager {
	return ListManager{ritHome: ritHome, file: file}
}

func (li ListManager) List() (formula.Repos, error) {
	repos := formula.Repos{}
	reposFilePath := path.Join(li.ritHome, reposDirName, reposFileName)
	if !li.file.Exists(reposFilePath) {
		return repos, nil
	}

	file, err := li.file.Read(reposFilePath)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(file, &repos); err != nil {
		return nil, err
	}

	sort.Sort(repos)

	return repos, nil
}
