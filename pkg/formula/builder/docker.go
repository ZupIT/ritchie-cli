/*
 * Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package builder

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"os/user"
	"runtime"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/os/osutil"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

const (
	defaultContainerCmd = "cd /app && /usr/bin/make build"
	volumePattern       = "%s:/app"
)

var ErrDockerBuild = errors.New("failed building formula with Docker, we will try to build your formula locally")

type DockerManager struct{}

func NewBuildDocker() formula.DockerBuilder {
	return DockerManager{}
}

func (do DockerManager) Build(formulaPath, dockerImg string) error {
	volume := fmt.Sprintf(volumePattern, formulaPath)
	containerCmd, err := containerCmd()
	if err != nil {
		return err
	}

	args := []string{"run", "-u", "0:0", "-v", volume, "--entrypoint", "/bin/sh", dockerImg, "-c", containerCmd}

	var stderr bytes.Buffer
	cmd := exec.Command("docker", args...)
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if stderr.Bytes() != nil {
			prompt.Error(stderr.String())
		}
		return ErrDockerBuild
	}

	return nil
}

func containerCmd() (string, error) {
	os := runtime.GOOS
	switch os {
	case osutil.Windows:
		return defaultContainerCmd, nil
	default:
		currentUser, err := user.Current()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s && chown -R %s bin", defaultContainerCmd, currentUser.Uid), nil
	}
}
