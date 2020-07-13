package repo

import "github.com/ZupIT/ritchie-cli/pkg/formula"

type ListUpdateManager struct {
	formula.RepositoryLister
	formula.RepositoryUpdater
}

func NewListUpdater(repoList formula.RepositoryLister, repoUpdate formula.RepositoryUpdater) formula.RepositoryListUpdater {
	return ListUpdateManager{
		RepositoryLister:  repoList,
		RepositoryUpdater: repoUpdate,
	}
}
