package formula

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
)

func (d DefaultRunner) PreRun(def Definition) (RunData, error) {
	pwd, _ := os.Getwd()
	formulaPath := def.FormulaPath(d.ritchieHome)
	config, err := d.loadConfig(formulaPath, def)
	if err != nil {
		return RunData{}, err
	}

	binName := def.BinName()
	binPath := def.BinPath(formulaPath)
	binFilePath := def.BinFilePath(binPath, binName)

	if !fileutil.Exists(binFilePath) {
		url := def.BundleUrl()
		name := def.BundleName()
		zipFile, err := d.downloadFormulaBundle(url, formulaPath, name)
		if err != nil {
			return RunData{}, err
		}

		if err := d.unzipFile(zipFile, formulaPath); err != nil {
			return RunData{}, err
		}
	}

	tmpDir, tmpBinDir, err := d.createWorkDir(binPath, def)
	if err != nil {
		return RunData{}, err
	}

	if err := os.Chdir(tmpBinDir); err != nil {
		return RunData{}, err
	}

	tmpBinFilePath := def.BinFilePath(tmpBinDir, binName)

	run := RunData{
		pwd:            pwd,
		formulaPath:    formulaPath,
		binPath:        binPath,
		tmpDir:         tmpDir,
		tmpBinDir:      tmpBinDir,
		tmpBinFilePath: tmpBinFilePath,
		config:         config,
	}

	return run, nil
}

func (d DefaultRunner) loadConfig(formulaPath string, def Definition) (Config, error) { // Pre run
	configName := def.ConfigName()
	configPath := def.ConfigPath(formulaPath, configName)
	if !fileutil.Exists(configPath) {
		url := def.ConfigUrl(configName)
		if err := d.downloadConfig(url, formulaPath, configName); err != nil {
			return Config{}, err
		}
	}

	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}

	var formulaConfig Config
	if err := json.Unmarshal(configFile, &formulaConfig); err != nil {
		return Config{}, err
	}
	return formulaConfig, nil
}

func (d DefaultRunner) createWorkDir(binPath string, def Definition) (string, string, error) {
	u := uuid.New().String()
	tDir, tBDir := def.TmpWorkDirPath(d.ritchieHome, u)

	if err := fileutil.CreateDirIfNotExists(tBDir, 0755); err != nil {
		return "", "", err
	}

	if err := fileutil.CopyDirectory(binPath, tBDir); err != nil {
		return "", "", err
	}

	return tDir, tBDir, nil
}

func (d DefaultRunner) downloadFormulaBundle(url, destPath, zipName string) (string, error) {
	log.Println("Download formula...")

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("formula bin not found")
	}

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusNotFound:
		return "", errors.New("formula bin not found")
	default:
		return "", errors.New("unknown error when downloading your formula")
	}

	file := fmt.Sprintf("%s/%s", destPath, zipName)

	if err := fileutil.CreateDirIfNotExists(destPath, 0755); err != nil {
		return "", err
	}
	out, err := os.Create(file)
	if err != nil {
		return "", err
	}

	defer out.Close()
	if _, err = io.Copy(out, resp.Body); err != nil {
		return "", err
	}

	log.Println("Done.")
	return file, nil
}

func (d DefaultRunner) downloadConfig(url, destPath, configName string) error {
	log.Println("Downloading config file...")

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusNotFound:
		return errors.New("config file not found")
	default:
		return errors.New("unknown error when downloading your config file")
	}

	file := fmt.Sprintf("%s/%s", destPath, configName)

	if err := fileutil.CreateDirIfNotExists(destPath, 0755); err != nil {
		return err
	}

	out, err := os.Create(file)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return err
	}

	log.Println("Done.")
	return nil
}

func (d DefaultRunner) unzipFile(filename, destPath string) error {
	log.Println("Installing formula...")

	if err := fileutil.CreateDirIfNotExists(destPath, 0655); err != nil {
		return err
	}

	if err := fileutil.Unzip(filename, destPath); err != nil {
		return err
	}

	if err := fileutil.RemoveFile(filename); err != nil {
		return err
	}

	log.Println("Done.")
	return nil
}
