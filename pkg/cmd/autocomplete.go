package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/autocomplete"
)

const (
	zsh        autocomplete.ShellName = "zsh"
	bash       autocomplete.ShellName = "bash"
	fish       autocomplete.ShellName = "fish"
	powerShell autocomplete.ShellName = "powershell"
)

// autocompleteCmd type for set autocomplete command
type autocompleteCmd struct {
	autocomplete.Generator
}

// NewAutocompleteCmd creates a new cmd instance
func NewAutocompleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "completion SUBCOMMAND",
		Short:   "Add autocomplete for terminal",
		Long:    `Add autocomplete for terminal, Available for (bash, zsh).`,
		Example: "rit completion zsh",
	}
}

// NewAutocompleteZsh creates a new cmd instance zsh
func NewAutocompleteZsh(g autocomplete.Generator) *cobra.Command {
	a := &autocompleteCmd{g}

	return &cobra.Command{
		Use:     zsh.String(),
		Short:   "Add zsh autocomplete for terminal",
		Long:    "Add zsh autocomplete for terminal",
		Example: "rit completion zsh",
		RunE:    a.runFunc(),
	}
}

// NewAutocompleteBash creates a new cmd instance zsh
func NewAutocompleteBash(g autocomplete.Generator) *cobra.Command {
	a := &autocompleteCmd{g}

	return &cobra.Command{
		Use:     bash.String(),
		Short:   "Add bash autocomplete for terminal",
		Long:    "Add bash autocomplete for terminal",
		Example: "rit completion bash",
		RunE:    a.runFunc(),
	}
}

// NewAutocompleteFish creates a new cmd instance fish
func NewAutocompleteFish(g autocomplete.Generator) *cobra.Command {
	a := &autocompleteCmd{g}

	return &cobra.Command{
		Use:     fish.String(),
		Short:   "Add fish(>3.X) autocomplete for terminal",
		Long:    "Add fish(>3.X) autocomplete for terminal, fish 2.0 is not supported only fish 3.X",
		Example: "rit completion fish",
		RunE:    a.runFunc(),
	}
}

// NewAutocompletePowerShell creates a new cmd instance PowerShell
func NewAutocompletePowerShell(g autocomplete.Generator) *cobra.Command {
	a := &autocompleteCmd{g}

	return &cobra.Command{
		Use:     powerShell.String(),
		Short:   "Add powerShell autocomplete for terminal",
		Long:    "Add powerShell autocomplete for terminal",
		Example: "rit completion powershell",
		RunE:    a.runFunc(),
	}
}

func (a autocompleteCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		s := autocomplete.ShellName(cmd.Use)
		c, err := a.Generate(s, cmd)
		if err != nil {
			return err
		}

		fmt.Println(c)
		return nil
	}

}
