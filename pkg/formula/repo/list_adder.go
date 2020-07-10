package repo

import "github.com/ZupIT/ritchie-cli/pkg/formula"

type ListAddManager struct {
	formula.RepositoryAdder
	formula.RepositoryLister
}

func NewListAdder(repoList formula.RepositoryLister, repoAdd formula.RepositoryAdder) formula.RepositoryAddLister {
	return ListAddManager{
		RepositoryLister: repoList,
		RepositoryAdder:  repoAdd,
	}
}
