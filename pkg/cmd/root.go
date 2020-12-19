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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
	"github.com/ZupIT/ritchie-cli/pkg/sliceutil"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/version"
)

const (
	latestVersionMsg            = "Latest available version: %s"
	versionMsg                  = "%s\n  Build date: %s\n  Built with: %s\n"
	versionMsgWithLatestVersion = "%s\n  %s\n  Build date: %s\n  Built with: %s\n"
	cmdUse                      = "rit"
	cmdShortDescription         = "rit is a NoOps CLI"
	cmdDescription              = `A CLI to create, store and share any kind of 
automations, executing them through command lines.
Complete documentation available at https://github.com/ZupIT/ritchie-cli`
)

var (
	Version   = ""
	BuildDate = "unknown"
	MsgInit   = "To start using rit, you need to initialize rit first.\nCommand: rit init"

	allowList = []string{
		cmdUse,
		fmt.Sprintf("%s help", cmdUse),
		fmt.Sprintf("%s completion zsh", cmdUse),
		fmt.Sprintf("%s completion bash", cmdUse),
		fmt.Sprintf("%s completion fish", cmdUse),
		fmt.Sprintf("%s completion powershell", cmdUse),
		fmt.Sprintf("%s init", cmdUse),
		fmt.Sprintf("%s upgrade", cmdUse),
		fmt.Sprintf("%s add repo", cmdUse),
	}
	upgradeList = []string{
		cmdUse,
	}

	blockedCmdsByCommons = []string{
		fmt.Sprintf("%s create formula", cmdUse),
	}
)

type rootCmd struct {
	ritchieHome string
	dir         stream.DirCreateChecker
	file        stream.FileWriteReadExistRemover
	tutorial    rtutorial.Finder
	version     version.Manager
	tree        formula.TreeGenerator
	repo        formula.RepositoryListWriter
}

func NewRootCmd(
	ritchieHome string,
	dir stream.DirCreateChecker,
	file stream.FileWriteReadExistRemover,
	tutorial rtutorial.Finder,
	version version.Manager,
	tree formula.TreeGenerator,
	repo formula.RepositoryListWriter,
) *cobra.Command {
	o := &rootCmd{
		ritchieHome: ritchieHome,
		dir:         dir,
		file:        file,
		tutorial:    tutorial,
		version:     version,
		tree:        tree,
		repo:        repo,
	}

	cmd := &cobra.Command{
		Use:                cmdUse,
		Short:              cmdShortDescription,
		Long:               cmdDescription,
		Version:            o.versionFlag(),
		PersistentPreRunE:  o.PreRunFunc(),
		PersistentPostRunE: o.PostRunFunc(),
		RunE:               runHelp,
		SilenceErrors:      true,
		TraverseChildren:   true,
		ValidArgs:          []string{""},
		Args:               cobra.OnlyValidArgs,
	}
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
	return cmd
}

func (ro *rootCmd) PreRunFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		if err := ro.dir.Create(ro.ritchieHome); err != nil {
			return err
		}

		if isUpgradeCommand(allowList, cmd) || isCompleteCmd(cmd) {
			return nil
		}

		if !ro.ritchieIsInitialized() && isBlockedByCommons(blockedCmdsByCommons, cmd) {
			fmt.Println(MsgInit)
			os.Exit(0)
		}

		if err := ro.convertTree(); err != nil {
			return err
		}

		if err := ro.convertContextsFileToEnvsFile(); err != nil {
			return err
		}

		return nil
	}
}

func (ro *rootCmd) PostRunFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		printNewVersionMessage(cmd, ro)
		if !ro.ritchieIsInitialized() && cmd.Use == cmdUse {
			tutorialHolder, err := ro.tutorial.Find()
			if err != nil {
				return err
			}
			tutorialRit(tutorialHolder.Current)
		}
		return nil
	}
}

func printNewVersionMessage(cmd *cobra.Command, ro *rootCmd) {
	if isUpgradeCommand(upgradeList, cmd) {
		currentStable, _ := ro.version.StableVersion()
		prompt.Warning(ro.version.VerifyNewVersion(currentStable, Version))
	}
}

func isUpgradeCommand(upgradeList []string, cmd *cobra.Command) bool {
	return sliceutil.Contains(upgradeList, cmd.CommandPath())
}

func isCompleteCmd(cmd *cobra.Command) bool {
	return strings.Contains(cmd.CommandPath(), "__complete")
}

func isBlockedByCommons(blockList []string, cmd *cobra.Command) bool {
	return sliceutil.Contains(blockList, cmd.CommandPath())
}

func (ro *rootCmd) versionFlag() string {
	latestVersion, err := ro.version.StableVersion()
	if err == nil && latestVersion != Version {
		formattedLatestVersionMsg := prompt.Yellow(fmt.Sprintf(latestVersionMsg, latestVersion))
		return fmt.Sprintf(
			versionMsgWithLatestVersion,
			Version,
			formattedLatestVersionMsg,
			BuildDate,
			runtime.Version())
	}
	return fmt.Sprintf(versionMsg, Version, BuildDate, runtime.Version())
}

func runHelp(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}

func tutorialRit(tutorialStatus string) {
	const tagTutorial = "\n[TUTORIAL]"
	const MessageTitle = "To initialize the ritchie:"
	const MessageBody = ` ∙ Run "rit init"` + "\n"

	if tutorialStatus == tutorialStatusEnabled {
		prompt.Info(tagTutorial)
		prompt.Info(MessageTitle)
		fmt.Print(MessageBody)
	}
}

// TODO: remove this method in the next release
func (ro *rootCmd) convertContextsFileToEnvsFile() error {
	ctx := struct {
		Current string   `json:"current_context"`
		All     []string `json:"contexts"`
	}{}

	contextsPath := filepath.Join(ro.ritchieHome, "contexts")
	if !ro.file.Exists(contextsPath) {
		return nil
	}

	bytes, err := ro.file.Read(contextsPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &ctx); err != nil {
		return err
	}

	envsPath := filepath.Join(ro.ritchieHome, env.FileName)
	envHolder := env.Holder{
		Current: ctx.Current,
		All:     ctx.All,
	}

	envs, err := json.Marshal(envHolder)
	if err != nil {
		return err
	}

	if err := ro.file.Write(envsPath, envs); err != nil {
		return err
	}

	if err := ro.file.Remove(contextsPath); err != nil {
		return err
	}

	return nil
}

func (ro *rootCmd) ritchieIsInitialized() bool {
	commonsRepoPath := filepath.Join(ro.ritchieHome, "repos", "commons")
	return ro.dir.Exists(commonsRepoPath)
}

// TODO: remove this method in the next release
func (ro *rootCmd) convertTree() error {
	repos, err := ro.repo.List()
	if err != nil {
		return err
	}

	var hasUpdate bool
	wg := sync.WaitGroup{}
	for i := range repos {
		if repos[i].TreeVersion == tree.Version {
			continue
		}

		wg.Add(1)
		go ro.generateTree(repos[i].Name.String(), &wg)

		repos[i].TreeVersion = tree.Version
		hasUpdate = true
	}
	wg.Wait()

	if !hasUpdate {
		return nil
	}

	if err := ro.repo.Write(repos); err != nil {
		return err
	}

	return nil
}

func (ro *rootCmd) generateTree(repo string, wg *sync.WaitGroup) {
	defer wg.Done()

	repoPath := filepath.Join(ro.ritchieHome, "repos", repo)
	tree, err := ro.tree.Generate(repoPath)
	if err != nil {
		return
	}

	bb, err := json.MarshalIndent(tree, "", "\t")
	if err != nil {
		return
	}

	treePath := filepath.Join(repoPath, "tree.json")
	if err := ro.file.Write(treePath, bb); err != nil {
		return
	}
}
