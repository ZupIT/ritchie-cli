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
	"path/filepath"
	"runtime"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/os/osutil"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const (
	containerCmdFmt   = "%s && chown -R %s bin"
	makeContainerCmd  = "cd /app && /usr/bin/make build"
	shellContainerCmd = "cd /app && ./build.sh"
	volumePattern     = "%s:/app"
)

var ErrDockerBuild = errors.New("failed building formula with Docker, we will try to build your formula locally")

type Manager struct {
	file stream.FileExister
}

func NewBuildDocker(file stream.FileExister) formula.DockerBuilder {
	return Manager{file: file}
}

func (do Manager) Build(formulaPath, dockerImg string) error {
	volume := fmt.Sprintf(volumePattern, formulaPath)
	containerCmd, err := do.containerCmd(formulaPath)
	if err != nil {
		return err
	}

	args := []string{"run", "--rm", "-u", "0:0", "-v", volume, "--entrypoint", "/bin/sh", dockerImg, "-c", containerCmd}

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

func (do Manager) containerCmd(formulaPath string) (string, error) {
	os := runtime.GOOS
	switch os {
	case osutil.Windows:
		return makeContainerCmd, nil
	default:
		currentUser, err := user.Current()
		if err != nil {
			return "", err
		}

		execFile := filepath.Join(formulaPath, buildSh)
		if do.file.Exists(execFile) {
			return fmt.Sprintf(containerCmdFmt, shellContainerCmd, currentUser.Uid), nil
		}

		return fmt.Sprintf(containerCmdFmt, makeContainerCmd, currentUser.Uid), nil
	}
}
