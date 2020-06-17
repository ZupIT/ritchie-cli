package formula

type DefaultPreRunner struct {
	sDefault Setuper
}

func NewDefaultPreRunner(setuper Setuper) DefaultPreRunner {
	return DefaultPreRunner{sDefault: setuper}
}

func (d DefaultPreRunner) PreRun(def Definition) (Setup, error) {
	return d.sDefault.Setup(def)
}
