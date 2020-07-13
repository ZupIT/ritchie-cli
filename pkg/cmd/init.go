package cmd

import (
	"time"

	"github.com/kaduartur/go-cli-spinner/pkg/spinner"
	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/github"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

const UsageMsg = ` How to contribute new formulas to the Ritchie community?
 You must fork the Github repository "https://github.com/ZupIT/ritchie-formulas" 
 and then follow the step by step below:
  ∙ git clone https://github.com/{{your_github_user}}/ritchie-formulas
  ∙ Run the command "rit create formula" and add the location where you cloned your 
    repository. Rit will create a formula template that you can already test.
  ∙ Open the project with your favorite text editor.
  ∙ In order to test your new formula, you can run the command "rit build formula" or
    "rit build formula --watch" to have automatic updates when editing your formula.`

var CommonsRepoURL = "https://github.com/kaduartur/ritchie-formulas"

type InitCmd struct {
	repo formula.RepositoryAdder
	git  github.Repositories
}

// NewInitCmd creates init command for single edition
func NewInitCmd(repo formula.RepositoryAdder, git github.Repositories) *cobra.Command {
	o := InitCmd{repo: repo, git: git}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Init rit",
		Long:  "Initialize rit configuration",
		RunE:  o.runPrompt(),
	}

	return cmd
}

func (in InitCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		repo := formula.Repo{
			Name:     "commons",
			Url:      CommonsRepoURL,
			Priority: 0,
		}

		s := spinner.StartNew("We are adding the commons repository, wait a moment, please...")
		time.Sleep(time.Second * 2)

		repoInfo := github.NewRepoInfo(repo.Url, repo.Token)

		tag, err := in.git.LatestTag(repoInfo)
		if err != nil {
			return err
		}

		repo.Version = formula.RepoVersion(tag.Name)

		if err := in.repo.Add(repo); err != nil {
			return err
		}

		s.Success(prompt.Green("Okay, now you can use rit.\n"))
		prompt.Info(UsageMsg)

		return nil
	}
}
