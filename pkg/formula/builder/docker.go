package builder

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"os/user"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

const volumePattern = "%s:/app"

var ErrDockerBuild = errors.New("failed building formula with Docker, we will try to build your formula locally")

type DockerManager struct{}

func NewBuildDocker() formula.DockerBuilder {
	return DockerManager{}
}

func (do DockerManager) Build(formulaPath, dockerImg string) error {
	volume := fmt.Sprintf(volumePattern, formulaPath)
	currentUser, err := user.Current()
	if err != nil {
		return err
	}

	containerCmd := fmt.Sprintf("cd /app && /usr/bin/make build && chown -R %s bin", currentUser.Uid)
	args := []string{"run", "-u", "0:0", "-v", volume, "--entrypoint", "/bin/sh", dockerImg, "-c", containerCmd}

	var stderr, stdout bytes.Buffer
	cmd := exec.Command("docker", args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return ErrDockerBuild
	}

	if stderr.String() != "" {
		prompt.Error("\n" + stderr.String())
		// return ErrDockerBuild
	}

	return nil
}
