package runner

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/session"
)

var (
	ErrInvalidRepoUrl            = prompt.NewError("RepoURL is invalid inside tree.json")
	ErrFormulaBinNotFound        = prompt.NewError("formula bin not found")
	ErrConfigFileNotFound        = prompt.NewError("config file not found")
	ErrUnknownFormulaDownload    = prompt.NewError("unknown error when downloading your formula")
	ErrUnknownConfigFileDownload = prompt.NewError("unknown error when downloading your config file")
	ErrCreateReqBundle           = prompt.NewError("failed to create request for bundle download")
	ErrCreateReqConfig           = prompt.NewError("failed to create request for config download")
)

type DefaultSetup struct {
	ritchieHome    string
	client         *http.Client
	sessionManager session.Manager
	edition        api.Edition
}

func NewDefaultSingleSetup(ritchieHome string, c *http.Client) DefaultSetup {
	return DefaultSetup{
		ritchieHome: ritchieHome,
		client:      c,
		edition:     api.Single,
	}
}

func NewDefaultTeamSetup(ritchieHome string, c *http.Client, sess session.Manager) DefaultSetup {
	return DefaultSetup{
		ritchieHome:    ritchieHome,
		client:         c,
		sessionManager: sess,
		edition:        api.Team,
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
			volume := fmt.Sprintf("%s:/app", formulaPath)
			args := []string{dockerRunCmd, "-v", volume, "--entrypoint", "/bin/sh", config.DockerIB, "-c",
				"cd /app/src && /usr/bin/make build"}
			cmd := exec.Command(docker, args...)
			cmd.Env = os.Environ()
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Start(); err != nil {
				prompt.Warning("Failed building formula with docker trying run local Makefile...")
				return buildMakefileLocal(fmt.Sprintf("%s/src", formulaPath))
			}

			if err := cmd.Wait(); err != nil {
				prompt.Warning("Failed building formula with docker trying run local Makefile...")
				return buildMakefileLocal(fmt.Sprintf("%s/src", formulaPath))
			}
			prompt.Info("\n\nSuccess building formula using docker...\n\n")
		} else {
			return buildMakefileLocal(fmt.Sprintf("%s/src", formulaPath))
		}
	}
	return nil
}

func buildMakefileLocal(formulaSrcPath string) error {
	prompt.Info("Building formula using local Makefile...")
	if err := os.Chdir(formulaSrcPath); err != nil {
		return err
	}
	cmd := exec.Command("make", "build")
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return errors.New("error build formula using make, verify your repository")
	}

	if err := cmd.Wait(); err != nil {
		return errors.New("error build formula using make, verify your repository")
	}

	prompt.Info("\n\nSuccess building formula using local Makefile...\n\n")
	return nil
}

func (d DefaultSetup) loadConfig(formulaPath string, def formula.Definition) (formula.Config, error) {
	configPath := def.ConfigPath(formulaPath)
	if !fileutil.Exists(configPath) {
		return formula.Config{}, fmt.Errorf("Failed load config file for formula."+
			"\nTry run rit update repo"+
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
