package formula

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
)

const docker = "docker"
const dockerBuildCmd = "build"
const dockerRunCmd = "run"
const dockerRemoveCmd = "rm"
const envFile = ".env"

type DockerRunner struct {
	PreRunner
	PostRunner
	InputRunner
}

func NewDockerRunner(preRunner PreRunner, postRunner PostRunner, inputRunner InputRunner) DockerRunner {
	return DockerRunner{preRunner, postRunner, inputRunner}
}

func (d DockerRunner) Run(def Definition, inputType api.TermInputType) error {
	setup, err := d.PreRun(def)
	if err != nil {
		return err
	}

	volume := fmt.Sprintf("%s:/app", setup.pwd)
	args := []string{dockerRunCmd, "--env-file", envFile, "-v", volume, "--name", setup.containerId, setup.containerId}
	cmd := exec.Command(docker, args...) // Run command "docker run -env-file .env -v "$(pwd):/app" --name (randomId) (randomId)"
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := d.Inputs(cmd, setup, inputType); err != nil {
		return err
	}

	for _, e := range cmd.Env { // Create a file named .env and add the environment variable inName=inValue
		if !fileutil.Exists(envFile) {
			if err := fileutil.WriteFile(envFile, []byte(e+ "\n")); err != nil {
				return err
			}
			continue
		}
		if err := fileutil.AppendFileData(envFile, []byte(e+"\n")); err != nil {
			return err
		}
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	if err := d.PostRun(setup, true); err != nil {
		return err
	}

	return nil
}
