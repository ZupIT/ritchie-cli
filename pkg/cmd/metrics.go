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
	metric.Checker
}

func NewMetricsCmd(file stream.FileWriteReadExister, inList prompt.InputList, checker metric.Checker) *cobra.Command {
	m := &metricsCmd{
		FileWriteReadExister: file,
		InputList:            inList,
		Checker: checker,
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
		path := metric.FilePath
		if !m.FileWriteReadExister.Exists(path) {
			options := []string{"yes", "no"}
			choose, err := m.InputList.List("You want to send anonymous data about the product, feature use, statistics and crash reports?", options)
			if err != nil {
				return err
			}

			err = m.FileWriteReadExister.Write(path, []byte(choose))
			if err != nil {
				return err
			}
			return nil
		}

		metricsStatus, err := m.Check()
		if err != nil {
			return err
		}

		changeTo := "yes"
		message := "You are now sending anonymous metrics. Thank you!"
		if metricsStatus {
			changeTo = "no"
			message = "You are no longer sending anonymous metrics."
		}

		err = m.FileWriteReadExister.Write(path, []byte(changeTo))
		if err != nil {
			return err
		}
		prompt.Info(message)
		return nil
	}
}
