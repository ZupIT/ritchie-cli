package repo

import "github.com/ZupIT/ritchie-cli/pkg/formula"

type ListWriteManager struct {
	formula.RepositoryLister
	formula.RepositoryWriter
}

func NewListWriter(
	repoList formula.RepositoryLister,
	repoWrite formula.RepositoryWriter,
) formula.RepositoryListWriter {
	return ListWriteManager{
		RepositoryLister: repoList,
		RepositoryWriter: repoWrite,
	}
}
