package runner

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/mattn/go-isatty"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
)

const (
	docker          = "docker"
	dockerBuildCmd  = "build"
	dockerRunCmd    = "run"
	dockerRemoveCmd = "rm"
	envFile         = ".env"
	isDocker        = true
)

type DockerRunner struct {
	formula.PreRunner
	formula.PostRunner
	formula.InputRunner
	ctxFinder rcontext.Finder
}

func NewDockerRunner(preRunner formula.PreRunner, postRunner formula.PostRunner, inputRunner formula.InputRunner, ctxFinder rcontext.Finder) DockerRunner {
	return DockerRunner{preRunner, postRunner, inputRunner, ctxFinder}
}

func (d DockerRunner) Run(def formula.Definition, inputType api.TermInputType, verboseFlag string) error {
	setup, err := d.PreRun(def)
	if err != nil {
		return err
	}

	volume := fmt.Sprintf("%s:/app", setup.Pwd)

	var args []string
	if isatty.IsTerminal(os.Stdout.Fd()) {
		args = []string{dockerRunCmd, "-it", "--env-file", envFile, "-v", volume, "--name", setup.ContainerId, setup.ContainerId}
	} else {
		args = []string{dockerRunCmd, "--env-file", envFile, "-v", volume, "--name", setup.ContainerId, setup.ContainerId}
	}

	cmd := exec.Command(docker, args...) // Run command "docker run -env-file .env -v "$(pwd):/app" --name (randomId) (randomId)"
	cmd.Env = os.Environ()

	verboseEnv := fmt.Sprintf(formula.EnvPattern, formula.VerboseEnv, verboseFlag)
	cmd.Env = append(cmd.Env, verboseEnv)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := d.Inputs(cmd, setup, inputType); err != nil {
		return err
	}

	for _, e := range cmd.Env { // Create a file named .env and add the environment variable inName=inValue
		if !fileutil.Exists(envFile) {
			if err := fileutil.WriteFile(envFile, []byte(e+"\n")); err != nil {
				return err
			}
			continue
		}
		if err := fileutil.AppendFileData(envFile, []byte(e+"\n")); err != nil {
			return err
		}
	}

	ctx, err := d.ctxFinder.Find()
	if err != nil {
		return err
	}

	if err := fileutil.AppendFileData(envFile, []byte("CONTEXT="+ctx.Current+"\n")); err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	if err := d.PostRun(setup, isDocker); err != nil {
		return err
	}

	return nil
}
