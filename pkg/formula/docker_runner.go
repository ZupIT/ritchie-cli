package formula

import (
	"fmt"
	"os"
	"os/exec"
)

const docker = "docker"
const dockerBuildCmd = "build"
const dockerRunCmd = "run"
const dockerRemoveCmd = "rm"
const envFile = ".env"

type DockerRunner struct {
	PreRunner
	InputRunner
}

func NewDockerRunner(preRunner PreRunner, inputRunner InputRunner) DockerRunner {
	return DockerRunner{preRunner, inputRunner}
}

func (d DockerRunner) Run(def Definition) error {
	setup, err := d.PreRun(def)
	if err != nil {
		return err
	}

	volume := fmt.Sprintf("%s:/app", setup.pwd)
	args := []string{dockerRunCmd, "-it", "--env-file", envFile, "-v", volume, "--name", setup.containerId, setup.containerId}
	cmd := exec.Command(docker, args...) // Run command "docker run -it -env-file .env -v "$(pwd):/app" --name (randomId) (randomId)"
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := d.Inputs(cmd, setup.formulaPath, &setup.config, true); err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	if err := PostRun(setup, true); err != nil {
		return err
	}

	return nil
}
