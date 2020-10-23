package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/metric"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type metricsCmd struct {
	file  stream.FileWriteReadExister
	input prompt.InputList
}

func NewMetricsCmd(file stream.FileWriteReadExister, inList prompt.InputList) *cobra.Command {
	m := &metricsCmd{
		file:  file,
		input: inList,
	}

	cmd := &cobra.Command{
		Use:       "metrics",
		Short:     "Turn metrics on and off",
		Long:      "Stop or start to send anonymous metrics to ritchie team.",
		RunE:      m.run(),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
	}

	return cmd

}

func (m metricsCmd) run() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		options := []string{"yes", "no"}
		choose, err := m.input.List("You want to send anonymous data "+
			"about the product, feature use, statistics and crash reports?",
			options)
		if err != nil {
			return err
		}

		message := "You are now sending anonymous metrics. Thank you!"
		if choose == "no" {
			message = "You are no longer sending anonymous metrics."
		}

		if err := m.file.Write(metric.FilePath, []byte(choose)); err != nil {
			return err
		}

		prompt.Info(message)
		return nil
	}
}
