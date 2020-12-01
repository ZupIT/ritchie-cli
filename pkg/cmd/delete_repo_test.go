package cmd

import "github.com/ZupIT/ritchie-cli/pkg/formula"

type repositoryDeleterMock struct {
	deleteMock func(repoName formula.RepoName) error
}

func (c repositoryDeleterMock) Delete(repoName formula.RepoName) error {
	return c.deleteMock(repoName)
}
