package formula

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/google/uuid"
)

type DockerPreRunner struct {
	sDefault Setuper
}

func NewDockerPreRunner(setuper Setuper) DockerPreRunner {
	return DockerPreRunner{sDefault: setuper}
}

func (d DockerPreRunner) PreRun(def Definition) (Setup, error) {
	setup, err := d.sDefault.Setup(def)
	if err != nil {
		return Setup{}, err
	}

	containerId, err := uuid.NewRandom()
	if err != nil {
		return Setup{}, err
	}

	setup.containerId = containerId.String()
	if err := buildImg(setup.containerId); err != nil {
		return Setup{}, err
	}

	return setup, nil
}

func buildImg(containerId string) error {
	fmt.Println("Building docker image...")
	args := []string{dockerBuildCmd, "-t", containerId, "."}
	cmd := exec.Command(docker, args...)
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	fmt.Println("Docker image was built :)")
	return nil
}
