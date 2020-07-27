package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/github"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

type updateRepoCmd struct {
	client *http.Client
	repo   formula.RepositoryListUpdater
	github github.Repositories
	prompt.InputText
	prompt.InputPassword
	prompt.InputURL
	prompt.InputList
	prompt.InputBool
	prompt.InputInt
}

func NewUpdateRepoCmd(
	client *http.Client,
	repo formula.RepositoryListUpdater,
	github github.Repositories,
	inText prompt.InputText,
	inPass prompt.InputPassword,
	inUrl prompt.InputURL,
	inList prompt.InputList,
	inBool prompt.InputBool,
	inInt prompt.InputInt,
) *cobra.Command {
	updateRepo := updateRepoCmd{
		client:        client,
		repo:          repo,
		github:        github,
		InputText:     inText,
		InputURL:      inUrl,
		InputList:     inList,
		InputBool:     inBool,
		InputInt:      inInt,
		InputPassword: inPass,
	}

	cmd := &cobra.Command{
		Use:     "repo",
		Short:   "Update a repository.",
		Example: "rit update repo",
		RunE:    RunFuncE(updateRepo.runStdin(), updateRepo.runPrompt()),
	}
	cmd.LocalFlags()

	return cmd
}

func (up updateRepoCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		repos, err := up.repo.List()
		if err != nil {
			return err
		}

		var reposName []string
		for i := range repos {
			reposName = append(reposName, repos[i].Name.String())
		}

		name, err := up.List("Select a repository to update: ", reposName)
		if err != nil {
			return err
		}

		var repo formula.Repo
		for i := range repos {
			if repos[i].Name == formula.RepoName(name) {
				repo = repos[i]
				break
			}
		}

		repoInfo := github.NewRepoInfo(repo.Url, repo.Token)
		tags, err := up.github.Tags(repoInfo)
		if err != nil {
			return err
		}

		version, err := up.List("Select your new version:", tags.Names())
		if err != nil {
			return err
		}

		if err := up.repo.Update(formula.RepoName(name), formula.RepoVersion(version)); err != nil {
			return err
		}

		successMsg := fmt.Sprintf("The %q repository was updated with success to version %q", name, version)
		prompt.Success(successMsg)

		return nil
	}
}

func (up updateRepoCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		r := formula.Repo{}

		err := stdin.ReadJson(os.Stdin, &r)
		if err != nil {
			return err
		}

		if err := up.repo.Update(r.Name, r.Version); err != nil {
			return err
		}

		successMsg := fmt.Sprintf("The %q repository was updated with success to version %q", r.Name, r.Version)
		prompt.Success(successMsg)

		return nil
	}
}
