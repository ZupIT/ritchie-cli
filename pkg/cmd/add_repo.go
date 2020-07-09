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

type AddRepoCmd struct {
	client *http.Client
	repo   formula.RepositoryAddLister
	github github.Repositories
	prompt.InputText
	prompt.InputPassword
	prompt.InputURL
	prompt.InputList
	prompt.InputBool
	prompt.InputInt
}

func NewAddRepoCmd(
	repo formula.RepositoryAddLister,
	github github.Repositories,
	inText prompt.InputText,
	inPass prompt.InputPassword,
	inUrl prompt.InputURL,
	inList prompt.InputList,
	inBool prompt.InputBool,
	inInt prompt.InputInt,
) *cobra.Command {
	addRepo := AddRepoCmd{
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
		Short:   "Add a repository.",
		Example: "rit add repo ",
		RunE:    RunFuncE(addRepo.runStdin(), addRepo.runPrompt()),
	}
	cmd.LocalFlags()

	return cmd
}

func (ad AddRepoCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		name, err := ad.Text("Repository name: ", true)
		if err != nil {
			return err
		}

		repos, err := ad.repo.List()
		if err != nil {
			return err
		}

		for i := range repos {
			repo := repos[i]
			if repo.Name == name {
				prompt.Warning(fmt.Sprintf("Your repository %q is gonna be overwritten.", repo.Name))
				choice, _ := ad.Bool("Want to proceed?", []string{"yes", "no"})
				if !choice {
					prompt.Info("Operation cancelled")
					return nil
				}
			}
		}

		isPrivate, err := ad.Bool("Is a private repository? ", []string{"no", "yes"})
		if err != nil {
			return err
		}

		var token string
		if isPrivate {
			token, err = ad.Password("Personal access tokens: ")
			if err != nil {
				return err
			}
		}

		url, err := ad.URL("Repository URL: ", "https://github.com/kaduartur/ritchie-formulas")
		if err != nil {
			return err
		}

		gitRepoInfo := github.NewRepoInfo(url, token)
		tags, err := ad.github.Tags(gitRepoInfo)
		if err != nil {
			return err
		}

		var tagNames []string
		for i := range tags {
			tagNames = append(tagNames, tags[i].Name)
		}

		version, err := ad.List("Select a tag version: ", tagNames)
		if err != nil {
			return err
		}

		priority, err := ad.Int("Set the priority [ps.: 0 is higher priority, the lower higher the priority] :")
		if err != nil {
			return err
		}

		repository := formula.Repo{
			Name:     name,
			Token:    token,
			Url:      url,
			Version:  version,
			Priority: int(priority),
		}

		if err := ad.repo.Add(repository); err != nil {
			return err
		}

		successMsg := fmt.Sprintf("The %q repository was added with success, now you can use your formulas with the Ritchie!", repository.Name)
		prompt.Success(successMsg)
		return nil
	}
}

func (ad AddRepoCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {

		r := formula.Repo{}

		err := stdin.ReadJson(os.Stdin, &r)
		if err != nil {
			prompt.Error(stdin.MsgInvalidInput)
			return err
		}

		if err := ad.repo.Add(r); err != nil {
			return err
		}

		successMsg := fmt.Sprintf("The %q repository was added with success, now you can use your formulas with the Ritchie!", r.Name)
		prompt.Success(successMsg)
		return nil
	}
}
