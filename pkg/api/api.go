package api

import (
	"fmt"
	"os/user"
	"strings"
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
		{Parent: "root_completion", Usage: "fish"},
		{Parent: "root_completion", Usage: "powershell"},
		{Parent: "root", Usage: "delete"},
		{Parent: "root_delete", Usage: "context"},
		{Parent: "root_delete", Usage: "repo"},
		{Parent: "root", Usage: "help"},
		{Parent: "root", Usage: "init"},
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
		{Parent: "root", Usage: "build"},
		{Parent: "root_build", Usage: "formula"},
		{Parent: "root", Usage: "upgrade"},
	}

	SingleCoreCmds = CoreCmds

	TeamCoreCmds = append(
		CoreCmds,
		[]Command{
			// temporarily removed {Parent: "root_create", Usage: "user"},
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
	Formula *Formula `json:"formula,omitempty"`
	Repo    string  `json:"Repo,omitempty"`
}

type Commands []Command

// Formula type
type Formula struct {
	Path    string `json:"path,omitempty"`
	Bin     string `json:"bin,omitempty"`
	LBin    string `json:"binLinux,omitempty"`
	MBin    string `json:"binDarwin,omitempty"`
	WBin    string `json:"binWindows,omitempty"`
	Bundle  string `json:"bundle,omitempty"`
	Config  string `json:"config,omitempty"`
	RepoURL string `json:"repoUrl,omitempty"`
}

// Edition type that represents Single or Team.
type Edition string

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

// UserHomeDir returns the home dir of the user
func UserHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return usr.HomeDir
}

// RitchieHomeDir returns the home dir of the ritchie
func RitchieHomeDir() string {
	return fmt.Sprintf(ritchieHomePattern, UserHomeDir())
}
