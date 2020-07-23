package cmd

import (
	"fmt"
	"time"

	"github.com/kaduartur/go-cli-spinner/pkg/spinner"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/github"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
)

var CommonsRepoURL = "https://github.com/zupIt/ritchie-formulas"

type initCmd struct {
	repo formula.RepositoryAdder
	git  github.Repositories
	rt   rtutorial.Finder
}

func NewInitCmd(repo formula.RepositoryAdder, git github.Repositories, rtf rtutorial.Finder) *cobra.Command {
	o := initCmd{repo: repo, git: git, rt: rtf}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Init rit",
		Long:  "Initialize rit configuration",
		RunE:  o.runPrompt(),
	}

	return cmd
}

func (in initCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		repo := formula.Repo{
			Name:     "commons",
			Url:      CommonsRepoURL,
			Priority: 0,
		}

		s := spinner.StartNew("Adding the commons repository...")
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

		tutorialHolder, err := in.rt.Find()
		if err != nil {
			return err
		}

		s.Success(prompt.Green("Initialization successful!"))

		tutorialInit(tutorialHolder.Current)
		return nil
	}
}

func tutorialInit(tutorialStatus string) {
	const tagTutorial = "[TUTORIAL]"
	const MessageTitle = ` How to create new formulas with Ritchie?`
	const MessageBody = ` ∙ Run "rit create formula"
 ∙ Open the project with your favorite text editor.\n
 `

	if tutorialStatus == tutorialStatusOn {
		prompt.Info("\n[TUTORIAL]")
		prompt.Info(MessageTitle)
		fmt.Println(MessageBody)
	}
}
