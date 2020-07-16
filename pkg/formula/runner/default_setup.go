package runner

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

var (
	ErrFormulaBinNotFound = prompt.NewError("formula bin not found")
	ErrConfigFileNotFound = prompt.NewError("config file not found")
)

const (
	makeCmd             = "make"
	buildCmd            = "build"
	volumeDockerPattern = "%s:/app"
)

type DefaultSetup struct {
	ritchieHome string
}

func NewDefaultSetup(ritchieHome string) DefaultSetup {
	return DefaultSetup{
		ritchieHome: ritchieHome,
	}
}

func (d DefaultSetup) Setup(def formula.Definition) (formula.Setup, error) {
	pwd, _ := os.Getwd()
	ritchieHome := d.ritchieHome
	formulaPath := def.FormulaPath(ritchieHome)
	config, err := d.loadConfig(formulaPath, def)
	if err != nil {
		return formula.Setup{}, err
	}

	binFilePath := def.BinFilePath(formulaPath)
	if err := d.buildFormula(formulaPath, binFilePath, config); err != nil {
		return formula.Setup{}, err
	}

	tmpDir, err := createWorkDir(ritchieHome, formulaPath, def)
	if err != nil {
		return formula.Setup{}, err
	}

	if err := os.Chdir(tmpDir); err != nil {
		return formula.Setup{}, err
	}

	s := formula.Setup{
		Pwd:         pwd,
		FormulaPath: formulaPath,
		BinName:     def.BinName(),
		BinPath:     def.BinPath(formulaPath),
		TmpDir:      tmpDir,
		Config:      config,
	}

	return s, nil
}

func (d DefaultSetup) buildFormula(formulaPath, binFilePath string, config formula.Config) error {
	if !fileutil.Exists(binFilePath) {
		if config.DockerIB != "" {
			prompt.Info("Building formula with docker...")
			volume := fmt.Sprintf(volumeDockerPattern, formulaPath)
			args := []string{dockerRunCmd, "-v", volume, "--entrypoint", "/bin/sh", config.DockerIB, "-c",
				"cd /app && /usr/bin/make build"}
			cmd := exec.Command(docker, args...)
			cmd.Env = os.Environ()
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Start(); err != nil {
				prompt.Warning("Failed building formula with docker, trying with local Makefile...")
				return buildMakefileLocal(formulaPath)
			}

			if err := cmd.Wait(); err != nil {
				prompt.Warning("Failed building formula with docker, trying with local Makefile...")
				return buildMakefileLocal(formulaPath)
			}
			prompt.Success("Successfully built formula using docker...")
		} else {
			return buildMakefileLocal(formulaPath)
		}
	}
	return nil
}

func buildMakefileLocal(formulaPath string) error {
	prompt.Info("Building formula using local Makefile...")
	if err := os.Chdir(formulaPath); err != nil {
		return err
	}
	cmd := exec.Command(makeCmd, buildCmd)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return errors.New("error building formula using make, verify your repository")
	}

	if err := cmd.Wait(); err != nil {
		return errors.New("error building formula using make, verify your repository")
	}

	prompt.Success("Successfully built formula using local Makefile...")
	return nil
}

func (d DefaultSetup) loadConfig(formulaPath string, def formula.Definition) (formula.Config, error) {
	configPath := def.ConfigPath(formulaPath)
	if !fileutil.Exists(configPath) {
		return formula.Config{}, fmt.Errorf("Load of config file for formula failed."+
			"\nTry running rit update repo"+
			"\nConfig file path not found: %s", configPath)
	}

	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return formula.Config{}, err
	}

	var formulaConfig formula.Config
	if err := json.Unmarshal(configFile, &formulaConfig); err != nil {
		return formula.Config{}, err
	}
	return formulaConfig, nil
}

func createWorkDir(home, formulaPath string, def formula.Definition) (string, error) {
	tDir := def.TmpWorkDirPath(home)

	if err := fileutil.CreateDirIfNotExists(tDir, 0755); err != nil {
		return "", err
	}

	if err := fileutil.CopyDirectory(def.BinPath(formulaPath), tDir); err != nil {
		return "", err
	}

	return tDir, nil
}
