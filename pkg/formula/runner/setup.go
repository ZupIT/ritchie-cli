package runner

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"time"

	"github.com/kaduartur/go-cli-spinner/pkg/spinner"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/os/osutil"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const loadConfigErrMsg = `Failed to load formula config file
Try running rit update repo
Config file path not found: %s`

type SetupManager struct {
	ritchieHome string
	make        formula.MakeBuilder
	docker      formula.DockerBuilder
	bat         formula.BatBuilder
	dir         stream.DirCreateListCopier
	file        stream.FileExister
}

func NewSetup(
	ritchieHome string,
	make formula.MakeBuilder,
	docker formula.DockerBuilder,
	bat formula.BatBuilder,
	dir stream.DirCreateListCopier,
	file stream.FileExister,
) SetupManager {
	return SetupManager{
		ritchieHome: ritchieHome,
		make:        make,
		docker:      docker,
		bat:         bat,
		dir:         dir,
		file:        file,
	}
}

func (d SetupManager) Setup(def formula.Definition, local bool) (formula.Setup, error) {
	pwd, _ := os.Getwd()
	formulaPath := def.FormulaPath(d.ritchieHome)
	config, err := d.loadConfig(formulaPath, def)
	if err != nil {
		return formula.Setup{}, err
	}

	binFilePath := def.BinFilePath(formulaPath)
	if !d.file.Exists(binFilePath) {
		s := spinner.StartNew("Building formula...")
		time.Sleep(2 * time.Second)
		if err := d.buildFormula(formulaPath, config.DockerIB, local); err != nil {
			s.Stop()
			return formula.Setup{}, err
		}
		s.Success(prompt.Green("Formula build with success!"))
	}

	tmpDir, err := d.createWorkDir(d.ritchieHome, formulaPath, def)
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

func (d SetupManager) buildFormula(formulaPath, dockerIB string, localFlag bool) error {
	if !localFlag && dockerIB != "" { // Build formula inside docker
		if err := d.docker.Build(formulaPath, dockerIB); err != nil {
			return err
		}
		return nil
	}

	if runtime.GOOS == osutil.Windows { // Build formula local with build.bat
		if err := d.bat.Build(formulaPath); err != nil {
			return err
		}
		return nil
	}

	if err := d.make.Build(formulaPath); err != nil { // Build formula local with Makefile
		return err
	}

	return nil
}

func (d SetupManager) loadConfig(formulaPath string, def formula.Definition) (formula.Config, error) {
	configPath := def.ConfigPath(formulaPath)
	if !d.file.Exists(configPath) {
		return formula.Config{}, fmt.Errorf(loadConfigErrMsg, configPath)
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

func (d SetupManager) createWorkDir(home, formulaPath string, def formula.Definition) (string, error) {
	tDir := def.TmpWorkDirPath(home)
	if err := d.dir.Create(tDir); err != nil {
		return "", err
	}

	if err := d.dir.Copy(def.BinPath(formulaPath), tDir); err != nil {
		return "", err
	}

	return tDir, nil
}
