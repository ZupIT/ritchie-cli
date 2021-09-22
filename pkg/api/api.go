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
	RootName        = "rit"
	ritchieHomeName = ".rit"
	CoreCmdsDesc    = "core commands:"
)

var (
	CoreCmds = Commands{
		"root_add":                   {Parent: "root", Usage: "add"},
		"root_add_repo":              {Parent: "root_add", Usage: "repo"},
		"root_add_workspace":         {Parent: "root_add", Usage: "workspace"},
		"root_completion":            {Parent: "root", Usage: "completion"},
		"root_completion_bash":       {Parent: "root_completion", Usage: "bash"},
		"root_completion_zsh":        {Parent: "root_completion", Usage: "zsh"},
		"root_completion_fish":       {Parent: "root_completion", Usage: "fish"},
		"root_completion_powershell": {Parent: "root_completion", Usage: "powershell"},
		"root_delete":                {Parent: "root", Usage: "delete"},
		"root_delete_env":            {Parent: "root_delete", Usage: "env"},
		"root_delete_repo":           {Parent: "root_delete", Usage: "repo"},
		"root_delete_workspace":      {Parent: "root_delete", Usage: "workspace"},
		"root_delete_formula":        {Parent: "root_delete", Usage: "formula"},
		"root_delete_credential":     {Parent: "root_delete", Usage: "credential"},
		"root_help":                  {Parent: "root", Usage: "help"},
		"root_init":                  {Parent: "root", Usage: "init"},
		"root_list":                  {Parent: "root", Usage: "list"},
		"root_list_repo":             {Parent: "root_list", Usage: "repo"},
		"root_list_credential":       {Parent: "root_list", Usage: "credential"},
		"root_list_formula":          {Parent: "root_list", Usage: "formula"},
		"root_list_workspace":        {Parent: "root_list", Usage: "workspace"},
		"root_set":                   {Parent: "root", Usage: "set"},
		"root_set_env":               {Parent: "root_set", Usage: "env"},
		"root_set_credential":        {Parent: "root_set", Usage: "credential"},
		"root_set_repo-priority":     {Parent: "root_set", Usage: "repo-priority"},
		"root_set_formula-runner":    {Parent: "root_set", Usage: "formula-runner"},
		"root_show":                  {Parent: "root", Usage: "show"},
		"root_show_env":              {Parent: "root_show", Usage: "env"},
		"root_show_formula-runner":   {Parent: "root_show", Usage: "formula-runner"},
		"root_create":                {Parent: "root", Usage: "create"},
		"root_create_formula":        {Parent: "root_create", Usage: "formula"},
		"root_update":                {Parent: "root", Usage: "update"},
		"root_update_repo":           {Parent: "root_update", Usage: "repo"},
		"root_update_workspace":      {Parent: "root_update", Usage: "workspace"},
		"root_build":                 {Parent: "root", Usage: "build"},
		"root_build_formula":         {Parent: "root_build", Usage: "formula"},
		"root_upgrade":               {Parent: "root", Usage: "upgrade"},
		"root_tutorial":              {Parent: "root", Usage: "tutorial"},
		"root_metrics":               {Parent: "root", Usage: "metrics"},
		"root_rename":                {Parent: "root", Usage: "rename"},
		"root_rename_formula":        {Parent: "root_rename", Usage: "formula"},
	}
)

type CommandID string

func (id CommandID) String() string {
	return string(id)
}

type ByLen []CommandID

func (a ByLen) Len() int {
	return len(a)
}

func (a ByLen) Less(i, j int) bool {
	return len(a[i]) < len(a[j])
}

func (a ByLen) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Command type
type Command struct {
	Parent         string `json:"parent"`
	Usage          string `json:"usage"`
	Help           string `json:"help,omitempty"`
	LongHelp       string `json:"longHelp,omitempty"`
	Formula        bool   `json:"formula,omitempty"`
	Repo           string `json:"-"`
	RepoNewVersion string `json:"-"`
}

type Commands map[CommandID]Command

// TermInputType represents the source of the inputs will be read
type TermInputType int

const (
	// Prompt input
	Prompt TermInputType = iota
	// Stdin input
	Stdin
	// Flag input
	Flag
)

func (t TermInputType) String() string {
	return [...]string{"Prompt", "Stdin", "Flag"}[t]
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
	newRitchieHomeName := os.Getenv("RITCHIE_HOME_DIR")
	if newRitchieHomeName != "" {
		return filepath.Join(UserHomeDir(), newRitchieHomeName)
	}
	return filepath.Join(UserHomeDir(), ritchieHomeName)
}
