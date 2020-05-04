package api

import (
	"fmt"
	"os/user"
)

const (
	ritchieHomePattern = "%s/.rit"
	// Team version
	Team = Edition("team")
	// Single version
	Single = Edition("single")
	// CoreCmdsDesc commands group description
	CoreCmdsDesc = "core commands:"
)

var (
	CoreCmds = []Command{
		{Parent: "root", Usage: "add"},
		{Parent: "root_add", Usage: "repo"},
		{Parent: "root", Usage: "completion"},
		{Parent: "root_completion", Usage: "bash"},
		{Parent: "root_completion", Usage: "zsh"},
		{Parent: "root", Usage: "delete"},
		{Parent: "root_delete", Usage: "context"},
		{Parent: "root_delete", Usage: "repo"},
		{Parent: "root", Usage: "help"},
		{Parent: "root", Usage: "list"},
		{Parent: "root_list", Usage: "repo"},
		{Parent: "root", Usage: "set"},
		{Parent: "root_set", Usage: "context"},
		{Parent: "root_set", Usage: "credential"},
		{Parent: "root", Usage: "show"},
		{Parent: "root_show", Usage: "context"},
		{Parent: "root", Usage: "create"},
		{Parent: "root_create", Usage: "formula"},
		{Parent: "root", Usage: "update"},
		{Parent: "root_update", Usage: "repo"},
		{Parent: "root", Usage: "clean"},
		{Parent: "root_clean", Usage: "repo"},
	}

	SingleCoreCmds = CoreCmds

	TeamCoreCmds = append(
		CoreCmds,
		[]Command{
			{Parent: "root_create", Usage: "user"},
			{Parent: "root_delete", Usage: "user"},
			{Parent: "root", Usage: "login"},
			{Parent: "root", Usage: "logout"},
		}...,
	)
)

// Command type
type Command struct {
	Parent  string  `json:"parent"`
	Usage   string  `json:"usage"`
	Help    string  `json:"help"`
	Formula Formula `json:"formula,omitempty"`
	Repo    string
}

// Formula type
type Formula struct {
	Path    string `json:"path"`
	Bin     string `json:"bin"`
	LBin    string `json:"binLinux"`
	MBin    string `json:"binDarwin"`
	WBin    string `json:"binWindows"`
	Bundle  string `json:"bundle"`
	Config  string `json:"config"`
	RepoURL string `json:"repoUrl"`
}

// Edition type that represents Single or Team.
type Edition string

func UserHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return usr.HomeDir
}

func RitchieHomeDir() string {
	return fmt.Sprintf(ritchieHomePattern, UserHomeDir())
}
