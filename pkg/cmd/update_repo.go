package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

type UpdateRepoCmd struct {
	client *http.Client
	repo   formula.RepositoryLister
	prompt.InputText
	prompt.InputPassword
	prompt.InputURL
	prompt.InputList
	prompt.InputBool
	prompt.InputInt
}

func NewUpdateRepoCmd(
	client *http.Client,
	repo formula.RepositoryLister,
	inText prompt.InputText,
	inPass prompt.InputPassword,
	inUrl prompt.InputURL,
	inList prompt.InputList,
	inBool prompt.InputBool,
	inInt prompt.InputInt,
) *cobra.Command {
	updateRepo := UpdateRepoCmd{
		client:        client,
		repo:          repo,
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

func (up UpdateRepoCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		repos, err := up.repo.List()
		if err != nil {
			return err
		}

		var reposName []string
		for i := range repos {
			reposName = append(reposName, repos[i].Name)
		}

		name, err := up.List("Select a repository to update: ", reposName)
		if err != nil {
			return err
		}

		var repo formula.Repo
		for i := range repos {
			if repos[i].Name == name {
				repo = repos[i]
				break
			}
		}

		fmt.Println(repo)
		return nil
	}
}

func (up UpdateRepoCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		return nil
	}
}
