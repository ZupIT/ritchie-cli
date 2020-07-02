package runner

import "github.com/ZupIT/ritchie-cli/pkg/formula"

type DefaultPreRunner struct {
	sDefault formula.Setuper
}

func NewDefaultPreRunner(setuper formula.Setuper) DefaultPreRunner {
	return DefaultPreRunner{sDefault: setuper}
}

func (d DefaultPreRunner) PreRun(def formula.Definition) (formula.Setup, error) {
	return d.sDefault.Setup(def)
}
