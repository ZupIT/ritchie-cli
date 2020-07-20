package builder

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

const (
	volumePattern              = "%s:/app"
	msgDockerBuildErrorPattern = `failed building formula with Docker, try run your formula with the flag "--local"
More about the error: %s`
)

type DockerManager struct{}

func NewBuildDocker() formula.DockerBuilder {
	return DockerManager{}
}

func (do DockerManager) Build(formulaPath, dockerImg string) error {
	volume := fmt.Sprintf(volumePattern, formulaPath)
	args := []string{"run", "-v", volume, "--entrypoint", "/bin/sh", dockerImg, "-c", "cd /app && /usr/bin/make build"}
	cmd := exec.Command("docker", args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if stderr.Bytes() != nil {
			errMsg := fmt.Sprintf(msgDockerBuildErrorPattern, stderr.String())
			return errors.New(errMsg)
		}
	}

	return nil
}
