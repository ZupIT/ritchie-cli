package cmd

import (
	"errors"
	"reflect"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/metric"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type metricsCmd struct {
	file  stream.FileWriteReadExister
	input prompt.InputList
}

var metricsFlagName = "metrics"
var metricsFlags = flags{
	{
		name:        metricsFlagName,
		shortName:   "",
		kind:        reflect.String,
		defValue:    "yes",
		description: "",
	},
}
var message string
var options = []string{"yes", "no"}

func NewMetricsCmd(file stream.FileWriteReadExister, inList prompt.InputList) *cobra.Command {
	m := &metricsCmd{
		file:  file,
		input: inList,
	}

	cmd := &cobra.Command{
		Use:       "metrics",
		Short:     "Turn metrics on and off",
		Long:      "Stop or start to send anonymous metrics to ritchie team.",
		RunE:      m.runCmd(),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
	}
	cmd.LocalFlags()
	addReservedFlags(cmd.Flags(), metricsFlags)
	return cmd

}

func (m metricsCmd) runCmd() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		choose, err := m.resolveInput(cmd)
		if err != nil {
			return err
		}
		if choose == "yes" {
			message = "You are now sending anonymous metrics. Thank you!"
		} else if choose == "no" {
			message = "You are no longer sending anonymous metrics."
		}

		if err := m.file.Write(metric.FilePath, []byte(choose)); err != nil {
			return err
		}

		prompt.Info(message)
		return nil
	}
}

func (m metricsCmd) resolveInput(cmd *cobra.Command) (string, error) {
	if IsFlagInput(cmd) {
		return m.runFlag(cmd)
	} else {
		return m.runPrompt()
	}
}

func (m metricsCmd) runFlag(cmd *cobra.Command) (string, error) {
	choose, err := cmd.Flags().GetString(metricsFlagName)
	if err != nil {
		return "", err
	}
	for i := range options {
		if strings.EqualFold(choose, options[i]) {
			break
		} else if i == len(options)-1{
			return "", errors.New("please provide a valid value to the flag metrics")
		}

	}
	return choose, nil
}

func (m metricsCmd) runPrompt() (string, error) {
	choose, err := m.input.List(
		"You want to send anonymous data about the product, feature use, statistics and crash reports?",
		options,
	)
	if err != nil {
		return "", err
	}
	return choose, nil
}
