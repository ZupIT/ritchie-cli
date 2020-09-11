/*
 * Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"fmt"
	"strings"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
)

const (
	totalReposMsg   = "There are %v repos"
	totalOneRepoMsg = "There is 1 repo"
)

type listRepoCmd struct {
	formula.RepositoryLister
	repoProviders formula.RepoProviders
	rt            rtutorial.Finder
}

func NewListRepoCmd(rl formula.RepositoryLister, rp formula.RepoProviders, rtf rtutorial.Finder) *cobra.Command {
	lr := listRepoCmd{rl, rp, rtf}
	cmd := &cobra.Command{
		Use:     "repo",
		Short:   "Show a list with all your available repositories",
		Example: "rit list repo",
		RunE:    lr.runFunc(),
	}
	return cmd
}

func (lr listRepoCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		repos, err := lr.List()
		if err != nil {
			return err
		}

		lr.printRepos(repos)

		if len(repos) != 1 {
			prompt.Info(fmt.Sprintf(totalReposMsg, len(repos)))
		} else {
			prompt.Info(totalOneRepoMsg)
		}

		tutorialHolder, err := lr.rt.Find()
		if err != nil {
			return err
		}
		tutorialListRepo(tutorialHolder.Current)
		return nil
	}
}

func (lr listRepoCmd) printRepos(repos formula.Repos) {
	maxTags := 5

	table := uitable.New()
	table.AddRow("PROVIDER", "NAME", "VERSION", "PRIORITY", "URL", "AVALIABLE VERSIONS")
	for _, repo := range repos {
		tagsOfRepo, _ := lr.getTagsAvaliable(repo)

		if len(tagsOfRepo) > maxTags {
			tagsOfRepo = tagsOfRepo[:maxTags]
			tagsOfRepo = append(tagsOfRepo, "...")
		}

		latestTags := strings.Join(tagsOfRepo, ", ")

		table.AddRow(repo.Provider, repo.Name, repo.Version, repo.Priority, repo.Url, latestTags)
	}
	raw := table.Bytes()
	raw = append(raw, []byte("\n")...)
	fmt.Println(string(raw))

}

func (lr listRepoCmd) getTagsAvaliable(repo formula.Repo) ([]string, error) {
	var listTagsStr []string
	formulaGit := lr.repoProviders.Resolve(repo.Provider)

	repoInfo := formulaGit.NewRepoInfo(repo.Url, repo.Token)
	tags, err := formulaGit.Repos.Tags(repoInfo)
	if err != nil {
		return []string{}, err
	}

	for _, tag := range tags {
		if tag.Name == repo.Version.String() {
			tag.Name = tag.Name + "(current)"
		}
		listTagsStr = append(listTagsStr, tag.Name)
	}

	return listTagsStr, nil
}

func tutorialListRepo(tutorialStatus string) {
	const tagTutorial = "\n[TUTORIAL]"
	const MessageTitle = "To update all repositories or delete an repository:"
	const MessageBody = ` ∙ Run "rit update repo"
 ∙ Run "rit delete repo"` + "\n"

	if tutorialStatus == tutorialStatusEnabled {
		prompt.Info(tagTutorial)
		prompt.Info(MessageTitle)
		fmt.Println(MessageBody)
	}
}
