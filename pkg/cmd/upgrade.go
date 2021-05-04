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
	"runtime"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/internal/pkg/i18n"
	"github.com/ZupIT/ritchie-cli/pkg/git"
	"github.com/ZupIT/ritchie-cli/pkg/git/github"
	"github.com/ZupIT/ritchie-cli/pkg/metric"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/upgrade"
	"github.com/ZupIT/ritchie-cli/pkg/version"
)

type UpgradeCmd struct {
	upgrade.Manager
	resolver version.Resolver
	upgrade.UrlFinder
	input  prompt.InputBool
	file   stream.FileWriteReadExister
	github git.Repositories
}

func NewUpgradeCmd(
	r version.Resolver,
	m upgrade.Manager,
	uf upgrade.UrlFinder,
	input prompt.InputBool,
	file stream.FileWriteReadExister,
	github git.Repositories,
) *cobra.Command {
	u := UpgradeCmd{
		Manager:   m,
		resolver:  r,
		UrlFinder: uf,
		input:     input,
		file:      file,
		github:    github,
	}

	return &cobra.Command{
		Use:       "upgrade",
		Short:     "Update rit version",
		Long:      `Update rit version to last stable version.`,
		RunE:      u.runFunc(),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
	}
}

func (u UpgradeCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		if !u.file.Exists(metric.FilePath) {
			header := i18n.T("init.metric.header")
			prompt.Info(header)
			fmt.Println(i18n.T("init.add.metric.info"))

			choose, err := u.input.Bool(AgreeSendMetrics, AcceptDeclineOpts)
			if err != nil {
				return err
			}

			responseToWrite := "yes"
			if !choose {
				responseToWrite = "no"
			}

			err = u.file.Write(metric.FilePath, []byte(responseToWrite))
			if err != nil {
				return err
			}
		}

		err := u.resolver.UpdateCache()
		if err != nil {
			return prompt.NewError(err.Error() + "\n")
		}

		upgradeURL := u.Url(runtime.GOOS)
		if err := u.Run(upgradeURL); err != nil {
			return prompt.NewError(err.Error() + "\n")
		}

		prompt.Success("Rit upgraded with success")
		repoInfo := github.NewRepoInfo(github.RitchieRepoURL, "")
		tag, err := u.github.LatestTag(repoInfo)
		if err != nil || tag.Description == "" {
			return nil
		}

		msg := fmt.Sprintf("Release %s description:", tag.Name)
		prompt.Info(msg)
		fmt.Println(tag.Description)

		return nil
	}
}
