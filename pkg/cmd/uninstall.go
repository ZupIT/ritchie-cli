package cmd

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/ZupIT/ritchie-cli/internal/pkg/config"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/spf13/cobra"
)

// uninstallCmd type for uninstall command.
type uninstallCmd struct {
	inBool        prompt.InputBool
	file          stream.FileRemover
	configManager config.Manager
}

// NewUninstallCmd creates a new cmd instance.
func NewUninstallCmd(
	inBool prompt.InputBool,
	file stream.FileRemover,
	configDeleter config.Manager,
) *cobra.Command {
	c := uninstallCmd{
		inBool:        inBool,
		file:          file,
		configManager: configDeleter,
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
		sure, err := c.inBool.Bool("This will remove rit from your computer, are you sure about that?", []string{"yes", "no"})
		if err != nil {
			fmt.Println(err)
		}

		if sure {
			if err := c.uninstall(runtime.GOOS); err != nil {
				fmt.Println(err)
			}
		}

		return nil
	}
}

func (c uninstallCmd) uninstall(os string) error {
	switch os {
	case "windows":
		fmt.Println("later")
	default:
		if err := c.removeBin(); err != nil {
			return err
		}
		if err := c.removeRitConfig(); err != nil {
			return err
		}
	}

	return nil
}

func (c uninstallCmd) removeBin() error {
	if err := c.file.Remove("/usr/local/bin/rit"); err != nil {
		return errors.New(
			"Fail to uninstall\n" +
				"Please try running this command again as root/Administrator\n" +
				"Example: sudo rit uninstall",
		)
	}
	return nil
}

func (c uninstallCmd) removeRitConfig() error {
	if err := c.configManager.Delete(); err != nil {
		return err
	}
	return nil
}
