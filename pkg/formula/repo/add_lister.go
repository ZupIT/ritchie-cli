package repo

import "github.com/ZupIT/ritchie-cli/pkg/formula"

type AddListManager struct {
	formula.RepositoryAdder
	formula.RepositoryLister
}

func NewAddLister(repoAdd formula.RepositoryAdder, repoList formula.RepositoryLister) formula.RepositoryAddLister {
	return AddListManager{
		RepositoryAdder:  repoAdd,
		RepositoryLister: repoList,
	}
}
