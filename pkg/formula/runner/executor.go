package runner

import (
	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

var _ formula.Executor = ExecutorManager{}

type ExecutorManager struct {
	runners formula.Runners
	config  formula.ConfigRunner
}

func NewExecutor(runners formula.Runners, config formula.ConfigRunner) ExecutorManager {
	return ExecutorManager{
		runners: runners,
		config:  config,
	}
}

func (ex ExecutorManager) Execute(exe formula.ExecuteData) error {
	runType := exe.RunType
	runner := ex.runners[runType]

	if runner == nil {
		configType, err := ex.config.Find()
		if err != nil {
			return err
		}

		runner = ex.runners[configType]
	}

	if err := runner.Run(exe.Def, exe.InType, exe.Verbose); err != nil {
		return err
	}

	return nil
}
