package repo

import "github.com/ZupIT/ritchie-cli/pkg/formula"

type ListCreateManager struct {
	formula.RepositoryLister
	formula.RepositoryCreator
}

func NewListCreator(repoList formula.RepositoryLister, repoCreate formula.RepositoryCreator) formula.RepositoryListCreator {
	return ListCreateManager{
		RepositoryLister:  repoList,
		RepositoryCreator: repoCreate,
	}
}
