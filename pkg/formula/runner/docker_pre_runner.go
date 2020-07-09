package runner

import (
	"fmt"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"os"
	"os/exec"

	"github.com/google/uuid"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

var ErrNotEnableDocker = prompt.NewError("this formula is not enabled to run in a container")
var ErrDockerNotFound =  prompt.NewError("you must have the docker installed on the machine to run formulas inside a container")

type DockerPreRunner struct {
	sDefault formula.Setuper
}

func NewDockerPreRunner(setuper formula.Setuper) DockerPreRunner {
	return DockerPreRunner{sDefault: setuper}
}

func (d DockerPreRunner) PreRun(def formula.Definition) (formula.Setup, error) {
	setup, err := d.sDefault.Setup(def) //TODO: Definir o setup para docker
	if err != nil {
		return formula.Setup{}, err
	}

	if err := validate(setup.TmpDir); err != nil {
		return formula.Setup{}, err
	}

	containerId, err := uuid.NewRandom()
	if err != nil {
		return formula.Setup{}, err
	}

	setup.ContainerId = containerId.String()
	if err := buildImg(setup.ContainerId); err != nil {
		return formula.Setup{}, err
	}

	return setup, nil
}

func validate(tmpBinDir string) error {
	args := []string{"version", "--format", "'{{.Server.Version}}'"}
	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	if output == nil || err != nil {
		return ErrDockerNotFound
	}

	dockerFile := fmt.Sprintf("%s/Dockerfile", tmpBinDir)
	if !fileutil.Exists(dockerFile) {
		return ErrNotEnableDocker
	}

	return nil
}

func buildImg(containerId string) error {
	fmt.Println("Building docker image...")
	args := []string{dockerBuildCmd, "-t", containerId, "."}
	cmd := exec.Command(docker, args...) // Run command "docker build -t (randomId) ."
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
