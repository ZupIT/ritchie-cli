package formula

import (
	"os/exec"

	"github.com/ZupIT/ritchie-cli/pkg/api"
)

var RunnerTypes = []string{"local", "docker"}

const (
	DefaultRun RunnerType = iota - 1
	LocalRun
	DockerRun
)

type RunnerType int

func (e RunnerType) Int() int {
	return int(e)
}

func (e RunnerType) String() string {
	return RunnerTypes[e]
}

type Runners map[RunnerType]Runner

type Executor interface {
	Execute(exe ExecuteData) error
}

type PreRunBuilder interface {
	Build(string)
}

type PreRunner interface {
	PreRun(def Definition) (Setup, error)
}

type Runner interface {
	Run(def Definition, inputType api.TermInputType, verbose bool) error
}

type PostRunner interface {
	PostRun(p Setup, docker bool) error
}

type InputRunner interface {
	Inputs(cmd *exec.Cmd, setup Setup, inputType api.TermInputType) error
}

type ConfigRunner interface {
	Create(runType RunnerType) error
	Find() (RunnerType, error)
}
