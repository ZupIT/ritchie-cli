package runner

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/google/uuid"
	"github.com/kaduartur/go-cli-spinner/pkg/spinner"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

const (
	dockerCmd = "docker"
	envFile   = ".env"
)

type RunManager struct {
	formula.PostRunner
	formula.InputRunner
	formula.Setuper
}

func NewFormulaRunner(postRun formula.PostRunner, input formula.InputRunner, setup formula.Setuper) formula.Runner {
	return RunManager{
		PostRunner:  postRun,
		InputRunner: input,
		Setuper:     setup,
	}
}

func (d RunManager) Run(def formula.Definition, inputType api.TermInputType, local bool) error {
	setup, err := d.Setup(def, local)
	if err != nil {
		return err
	}

	var isDocker bool
	var cmd *exec.Cmd
	if local || !validateDocker(setup.TmpDir) {
		cmd, err = d.RunLocal(setup, inputType)
		if err != nil {
			return err
		}
	} else {
		containerId, err := buildImg()
		if err != nil {
			return err
		}

		setup.ContainerId = containerId
		cmd, err = d.RunDocker(setup, inputType)
		if err != nil {
			return err
		}

		isDocker = true
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	if err := d.PostRun(setup, isDocker); err != nil {
		return err
	}

	return nil
}

func (d RunManager) RunDocker(setup formula.Setup, inputType api.TermInputType) (*exec.Cmd, error) {
	volume := fmt.Sprintf("%s:/app", setup.Pwd)
	args := []string{"run", "--env-file", envFile, "-v", volume, "--name", setup.ContainerId, setup.ContainerId}
	cmd := exec.Command(dockerCmd, args...) // Run command "docker run -env-file .env -v "$(pwd):/app" --name (randomId) (randomId)"
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := d.Inputs(cmd, setup, inputType); err != nil {
		return nil, err
	}

	for _, e := range cmd.Env { // Create a file named .env and add the environment variable inName=inValue
		if !fileutil.Exists(envFile) {
			if err := fileutil.WriteFile(envFile, []byte(e+"\n")); err != nil {
				return nil, err
			}
			continue
		}
		if err := fileutil.AppendFileData(envFile, []byte(e+"\n")); err != nil {
			return nil, err
		}
	}

	return cmd, nil
}

func (d RunManager) RunLocal(setup formula.Setup, inputType api.TermInputType) (*exec.Cmd, error) {
	formulaRun := path.Join(setup.TmpDir, setup.BinName)
	cmd := exec.Command(formulaRun)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = os.Environ()
	pwdEnv := fmt.Sprintf(formula.EnvPattern, formula.PwdEnv, setup.Pwd)
	cPwdEnv := fmt.Sprintf(formula.EnvPattern, formula.CPwdEnv, setup.Pwd)
	cmd.Env = append(cmd.Env, pwdEnv)
	cmd.Env = append(cmd.Env, cPwdEnv)

	if err := d.Inputs(cmd, setup, inputType); err != nil {
		return nil, err
	}

	return cmd, nil
}

func buildImg() (string, error) {
	s := spinner.StartNew("Building docker image to run...")
	containerId, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	args := []string{"build", "-t", containerId.String(), "."}
	cmd := exec.Command(dockerCmd, args...) // Run command "docker build -t (randomId) ."
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		s.Stop()
		return "", err
	}

	s.Success(prompt.Green("Docker image was built!"))
	return containerId.String(), err
}

// validate checks if able to run inside docker
func validateDocker(tmpDir string) bool {
	args := []string{"version", "--format", "'{{.Server.Version}}'"}
	cmd := exec.Command(dockerCmd, args...)
	output, err := cmd.CombinedOutput()
	if output == nil || err != nil {
		return false
	}

	dockerFile := path.Join(tmpDir, "Dockerfile")
	return fileutil.Exists(dockerFile)
}
