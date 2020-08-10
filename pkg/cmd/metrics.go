package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/metric"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type metricsCmd struct {
	stream.FileWriteReadExister
	prompt.InputList
}

func NewMetricsCmd(file stream.FileWriteReadExister, inList prompt.InputList) *cobra.Command {
	m := &metricsCmd{
		FileWriteReadExister: file,
		InputList:            inList,
	}

	cmd := &cobra.Command{
		Use:   "metrics",
		Short: "Turn metrics on and off",
		Long:  "Stop or start to send anonymous metrics to ritchie team.",
		RunE:  m.run(),
	}

	return cmd

}

func (m metricsCmd) run() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		if !m.FileWriteReadExister.Exists(metric.MetricsPath()) {
			options := []string{"yes", "no"}
			choose, err := m.InputList.List("You want to to send anonymous data about the product, feature use, statistics and crash reports?", options)
			if err != nil {
				return err
			}

			err = m.FileWriteReadExister.Write(metric.MetricsPath(), []byte(choose))
			if err != nil {
				return err
			}
			return nil
		}

		metricsStatus, err := m.FileWriteReadExister.Read(metric.MetricsPath())
		if err != nil {
			return err
		}

		changeTo := "no"
		message := "You are no longer sending anonymous metrics."
		if string(metricsStatus) == changeTo {
			changeTo = "yes"
			message = "You are now sending anonymous metrics. Thank you!"
		}

		err = m.FileWriteReadExister.Write(metric.MetricsPath(), []byte(changeTo))
		if err != nil {
			return err
		}
		prompt.Info(message)
		return nil
	}
}
