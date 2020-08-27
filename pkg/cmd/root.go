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
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
	"github.com/ZupIT/ritchie-cli/pkg/slice/sliceutil"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/version"
)

const (
	latestVersionMsg            = "Latest available version: %s"
	versionMsg                  = "%s\n  Build date: %s\n  Built with: %s\n"
	versionMsgWithLatestVersion = "%s\n  %s\n  Build date: %s\n  Built with: %s\n"
	cmdUse                      = "rit"
	cmdShortDescription         = "rit is a NoOps CLI"
	cmdDescription              = `A CLI that developers can build and operate
your applications without help from the infra staff.
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
)

type rootCmd struct {
	ritchieHome string
	dir         stream.DirCreateChecker
	rt          rtutorial.Finder
	vm          version.Manager
}

func NewRootCmd(
	ritchieHome string,
	dir stream.DirCreateChecker,
	rtf rtutorial.Finder,
	vm version.Manager,
) *cobra.Command { o := &rootCmd{
		ritchieHome: ritchieHome,
		dir:         dir,
		rt:          rtf,
		vm:          vm,
	}

	cmd := &cobra.Command{
		Use:                cmdUse,
		Short:              cmdShortDescription,
		Long:               cmdDescription,
		Version:            versionFlag(vm),
		PersistentPreRunE:  o.PreRunFunc(),
		PersistentPostRunE: o.PostRunFunc(),
		RunE:               runHelp,
		SilenceErrors:      true,
		TraverseChildren:   true,
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

		if !ro.ritchieIsInitialized() {
			fmt.Println(MsgInit)
			os.Exit(0)
		}
		return nil
	}
}

func (ro *rootCmd) PostRunFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		printNewVersionMessage(cmd, ro)
		if !ro.ritchieIsInitialized() && cmd.Use == cmdUse {
			tutorialHolder, err := ro.rt.Find()
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
		currentStable, _ := ro.vm.StableVersion()
		prompt.Warning(ro.vm.VerifyNewVersion(currentStable, Version))
	}
}

func isUpgradeCommand(upgradeList []string, cmd *cobra.Command) bool {
	return sliceutil.Contains(upgradeList, cmd.CommandPath())
}

func isCompleteCmd(cmd *cobra.Command) bool {
	return strings.Contains(cmd.CommandPath(), "__complete")
}

func versionFlag(vm version.Manager) string {
	latestVersion, err := vm.StableVersion()
	if err == nil && latestVersion != Version {
		formattedLatestVersionMsg := prompt.Yellow(fmt.Sprintf(latestVersionMsg, latestVersion))
		return fmt.Sprintf(versionMsgWithLatestVersion,
			Version,
			formattedLatestVersionMsg,
			BuildDate,
			runtime.Version())
	}
	return fmt.Sprintf(versionMsg,
		Version,
		BuildDate,
		runtime.Version())
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

func (ro *rootCmd) ritchieIsInitialized() bool {
	commonsRepoPath := filepath.Join(ro.ritchieHome, "repos", "commons")
	return ro.dir.Exists(commonsRepoPath)
}
