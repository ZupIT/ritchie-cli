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

package api

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

const (
	ritchieHomeName = ".rit"
	CoreCmdsDesc    = "core commands:"
)

var (
	CoreCmds = []Command{
		{Parent: "root", Usage: "add"},
		{Parent: "root_add", Usage: "repo"},
		{Parent: "root", Usage: "completion"},
		{Parent: "root_completion", Usage: "bash"},
		{Parent: "root_completion", Usage: "zsh"},
		{Parent: "root_completion", Usage: "fish"},
		{Parent: "root_completion", Usage: "powershell"},
		{Parent: "root", Usage: "delete"},
		{Parent: "root_delete", Usage: "context"},
		{Parent: "root_delete", Usage: "repo"},
		{Parent: "root_delete", Usage: "workspace"},
		{Parent: "root_delete", Usage: "formula"},
		{Parent: "root", Usage: "help"},
		{Parent: "root", Usage: "init"},
		{Parent: "root", Usage: "list"},
		{Parent: "root_list", Usage: "repo"},
		{Parent: "root_list", Usage: "credential"},
		{Parent: "root_list", Usage: "workspace"},
		{Parent: "root", Usage: "set"},
		{Parent: "root_set", Usage: "context"},
		{Parent: "root_set", Usage: "credential"},
		{Parent: "root_set", Usage: "repo-priority"},
		{Parent: "root_set", Usage: "formula-runner"},
		{Parent: "root", Usage: "show"},
		{Parent: "root_show", Usage: "context"},
		{Parent: "root_show", Usage: "formula-runner"},
		{Parent: "root", Usage: "create"},
		{Parent: "root_create", Usage: "formula"},
		{Parent: "root", Usage: "update"},
		{Parent: "root_update", Usage: "repo"},
		{Parent: "root", Usage: "build"},
		{Parent: "root_build", Usage: "formula"},
		{Parent: "root", Usage: "upgrade"},
		{Parent: "root", Usage: "tutorial"},
		{Parent: "root", Usage: "metrics"},
	}
)

// Command type
type Command struct {
	Id       string `json:"id"`
	Parent   string `json:"parent"`
	Usage    string `json:"usage"`
	Help     string `json:"help"`
	LongHelp string `json:"longHelp"`
	Formula  bool   `json:"formula,omitempty"`
	Repo     string `json:"Repo,omitempty"`
}

type Commands []Command

// TermInputType represents the source of the inputs will be readed
type TermInputType int

const (
	// Prompt input
	Prompt TermInputType = iota
	// Stdin input
	Stdin
)

func (t TermInputType) String() string {
	return [...]string{"Prompt", "Stdin"}[t]
}

// ToLower converts the input type to lower case
func (t TermInputType) ToLower() string {
	return strings.ToLower(t.String())
}

// UserHomeDir returns the home dir of the user,
// if rit is called with sudo, it returns the same path
func UserHomeDir() string {
	if os.Geteuid() == 0 {
		username := os.Getenv("SUDO_USER")
		if username != "" {
			if u, err := user.Lookup(username); err == nil {
				return u.HomeDir
			}
		}
	}

	usr, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return usr
}

// RitchieHomeDir returns the home dir of the ritchie
func RitchieHomeDir() string {
	return filepath.Join(UserHomeDir(), ritchieHomeName)
}
