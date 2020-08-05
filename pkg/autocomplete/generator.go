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

package autocomplete

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/slice/sliceutil"
)

const (
	binaryName              = "rit"
	binaryNamePattern       = "{{BinaryName}}"
	dynamicCodePattern      = "{{DynamicCode}}"
	autocompleteBashPattern = "{{AutoCompleteBash}}"
	rootCommandPattern      = "{{RootCommand}}"
	lastCommandPattern      = "{{LastCommand}}"
	funcNamePattern         = "{{FunctionName}}"
	commandsPattern         = "{{Commands}}"
	lineCommand             = "    commands+=(\"${command}\")"
	firstLevel              = "root"
	bashPattern             = "%s\n%s"
	bash                    = "bash"
	zsh                     = "zsh"
	fish                    = "fish"
	powerShell              = "powershell"
)

var (
	supportedAutocomplete = []string{bash, zsh, fish, powerShell}
	ErrNotSupported       = prompt.NewError("autocomplete for this terminal is not supported")
)

type GeneratorManager struct {
	treeManager formula.TreeManager
}

func NewGenerator(tm formula.TreeManager) GeneratorManager {
	return GeneratorManager{tm}
}

func (d GeneratorManager) Generate(s ShellName, cmd *cobra.Command) (string, error) {
	if !sliceutil.Contains(supportedAutocomplete, s.String()) {
		return "", ErrNotSupported
	}
	t := d.treeManager.MergedTree(true)

	var autocomplete string
	var err error
	switch s.String() {
	case bash:
		autocomplete = fmt.Sprintf(bashPattern, "#!/bin/bash", loadToBash(t))
	case zsh:
		autocomplete = loadToZsh(t)
	case fish:
		autocomplete, err = loadToFish(cmd)
	case powerShell:
		autocomplete, err = loadToPowerShell(cmd)
	}
	return autocomplete, err
}

func loadToPowerShell(cmd *cobra.Command) (string, error) {
	buffer := bytes.Buffer{}
	if err := cmd.Root().GenPowerShellCompletion(&buffer); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func loadToFish(cmd *cobra.Command) (string, error) {
	buffer := bytes.Buffer{}
	if err := cmd.Root().GenFishCompletion(&buffer, true); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func loadToBash(t formula.Tree) string {
	a := autoCompletionBash
	a = strings.Replace(a, binaryNamePattern, binaryName, -1)
	a = strings.Replace(a, dynamicCodePattern, loadDynamicCommands(t), 1)
	return a
}

func loadToZsh(t formula.Tree) string {
	a := autoCompletionZsh
	a = strings.Replace(a, binaryNamePattern, binaryName, -1)
	a = strings.Replace(a, autocompleteBashPattern, loadToBash(t), 1)
	return a
}

func loadDynamicCommands(t formula.Tree) string {
	mapCommand := loadCommands(t.Commands)
	bashCommands := loadBashCommands(mapCommand)

	var allCommands string
	for _, b := range bashCommands {
		functionName := formatterFunctionName(b.RootCommand)
		c := strings.Replace(command, rootCommandPattern, b.RootCommand, -1)
		c = strings.Replace(c, lastCommandPattern, b.LastCommand, -1)
		c = strings.Replace(c, funcNamePattern, functionName, -1)
		allCommands += strings.Replace(c, commandsPattern, b.Commands, -1)
	}
	return allCommands
}

func formatterFunctionName(funcName string) string {
	ff := strings.Split(funcName, "_")
	if len(ff) > 2 {
		funcName = ff[len(ff)-2] + "_" + ff[len(ff)-1]
	}
	return funcName
}

func loadCommands(cc []api.Command) map[string]CompletionCommand {
	tmpCmd := make(map[string]CompletionCommand)
	for _, v := range cc {
		c := tmpCmd[v.Parent]
		c.Content = append(c.Content, v.Usage)
		c.Before = v.Parent
		tmpCmd[v.Parent] = c
	}

	commands := make(map[string]CompletionCommand)
	for key, val := range tmpCmd {
		commands[key] = val
		for _, v := range val.Content {
			newKey := key + "_" + v
			if _, ok := tmpCmd[newKey]; !ok {
				commands[newKey] = CompletionCommand{
					Content: nil,
					Before:  newKey,
				}
			}
		}
	}
	return commands
}

func loadBashCommands(cc map[string]CompletionCommand) []BashCommand {
	var bb []BashCommand
	for key, val := range cc {
		rootCommand := key
		level := len(strings.Split(key, "_"))
		var commands string
		for _, v := range val.Content {
			commands += strings.Replace(lineCommand, "${command}", v, -1) + "\n"
		}
		if rootCommand == firstLevel {
			rootCommand = fmt.Sprintf("%s_%s", binaryName, rootCommand)
		}
		bb = append(bb, BashCommand{
			RootCommand: rootCommand,
			Commands:    commands,
			LastCommand: loadLastCommand(key),
			Level:       level,
		},
		)
	}
	return bb
}

func loadLastCommand(rootCommand string) string {
	cc := strings.Split(rootCommand, "_")
	return cc[len(cc)-1]
}
