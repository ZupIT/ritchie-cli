package repo

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type ListManager struct {
	ritHome string
	file    stream.FileReadExister
}

// ByPriority implements sort.Interface for []Repository based on
// the Priority field.
type ByPriority []formula.Repo

func (a ByPriority) Len() int           { return len(a) }
func (a ByPriority) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPriority) Less(i, j int) bool { return a[i].Priority < a[j].Priority }

func NewLister(ritHome string, file stream.FileReadExister) ListManager {
	return ListManager{
		ritHome: ritHome,
		file:    file,
	}
}

func (lm ListManager) List() ([]formula.Repo, error) {
	f, err := lm.loadReposFromDisk()
	if err != nil {
		return []formula.Repo{}, err
	}

	if len(f.Values) == 0 {
		return []formula.Repo{}, nil
	}

	sort.Sort(ByPriority(f.Values))

	return f.Values, nil
}

func (lm ListManager) loadReposFromDisk() (formula.RepoFile, error) {
	path := fmt.Sprintf(repositoryConfFilePattern, lm.ritHome)
	rf := formula.RepoFile{}

	if !lm.file.Exists(path) {
		return rf, nil
	}

	b, err := lm.file.Read(path)
	if err != nil {
		return rf, err
	}

	err = json.Unmarshal(b, &rf)
	return rf, err
}
