package cmd

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/spf13/cobra"
)

var (
	ErrCannotUninstall = errors.New("cannot remove rit")
)

// uninstallCmd type for uninstall command.
type uninstallCmd struct {
	inBool prompt.InputBool
	file   stream.FileRemover
}

// NewUninstallCmd creates a new cmd instance.
func NewUninstallCmd(
	inBool prompt.InputBool,
	file stream.FileRemover,
) *cobra.Command {
	c := uninstallCmd{
		inBool: inBool,
		file:   file,
	}

	cmd := &cobra.Command{
		Use:       "uninstall",
		Short:     "Uninstall rit",
		Example:   "rit uninstall",
		RunE:      c.runPrompt(),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
	}

	cmd.LocalFlags()

	return cmd
}

func (c uninstallCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		sure, err := c.inBool.Bool("Are you sure", []string{"yes", "no"})
		if err != nil {
			fmt.Println(err)
		}

		if sure {
			if err := c.uninstall(); err != nil {
				fmt.Println(err)
			}
		}

		return nil
	}
}

func (c uninstallCmd) uninstall() error {
	switch runtime.GOOS {
	case "windows":
		fmt.Println("later")
	default:
		if err := c.file.Remove("/usr/local/bin/rit"); err != nil {
			return errors.New(
				"Fail to uninstall\n" +
					"Please try running this command again as root/Administrator\n" +
					"Example: sudo rit uninstall",
			)
		}
	}

	return nil
}
