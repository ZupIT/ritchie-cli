package formula

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/docker/docker/pkg/urlutil"
	"github.com/google/uuid"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/http/headers"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/session"
)

var (
	ErrFormulaBinNotFound        = errors.New("formula bin not found")
	ErrConfigFileNotFound        = errors.New("config file not found")
	ErrUnknownFormulaDownload    = errors.New("unknown error when downloading your formula")
	ErrUnknownConfigFileDownload = errors.New("unknown error when downloading your config file")
	ErrCreateReqBundle           = errors.New("failed to create request for bundle download")
	ErrCreateReqConfig           = errors.New("failed to create request for config download")
	ErrInvalidRepoUrl            = errors.New("RepoURL is invalid inside tree.json")
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

func (d DefaultSetup) Setup(def Definition) (Setup, error) {
	pwd, _ := os.Getwd()
	ritchieHome := d.ritchieHome
	formulaPath := def.FormulaPath(ritchieHome)
	config, err := d.loadConfig(formulaPath, def)
	if err != nil {
		return Setup{}, err
	}

	binName := def.BinName()
	binPath := def.BinPath(formulaPath)
	binFilePath := def.BinFilePath(binPath, binName)
	if err := d.loadBundle(formulaPath, binFilePath, def); err != nil {
		return Setup{}, err
	}

	tmpDir, tmpBinDir, err := createWorkDir(ritchieHome, binPath, def)
	if err != nil {
		return Setup{}, err
	}

	if err := os.Chdir(tmpBinDir); err != nil {
		return Setup{}, err
	}

	tmpBinFilePath := def.BinFilePath(tmpBinDir, binName)

	s := Setup{
		pwd:            pwd,
		formulaPath:    formulaPath,
		binPath:        binPath,
		tmpDir:         tmpDir,
		tmpBinDir:      tmpBinDir,
		tmpBinFilePath: tmpBinFilePath,
		config:         config,
	}

	return s, nil
}

func (d DefaultSetup) loadConfig(formulaPath string, def Definition) (Config, error) {
	configName := def.ConfigName()
	configPath := def.ConfigPath(formulaPath, configName)
	if !fileutil.Exists(configPath) {
		url := def.ConfigURL(configName)
		if !urlutil.IsURL(url) {
			return Config{}, ErrInvalidRepoUrl
		}

		prompt.Info("Downloading formula config...")
		if err := d.downloadConfig(url, formulaPath, configName, def.RepoName); err != nil {
			return Config{}, err
		}
		prompt.Success("Formula config download completed!")
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

func (d DefaultSetup) loadBundle(formulaPath, binFilePath string, def Definition) error {
	if !fileutil.Exists(binFilePath) {
		url := def.BundleURL()
		if !urlutil.IsURL(url) {
			return ErrInvalidRepoUrl
		}

		name := def.BundleName()
		zipFile, err := d.downloadFormulaBundle(url, formulaPath, name, def.RepoName)
		if err != nil {
			return err
		}

		if err := unzipFile(zipFile, formulaPath); err != nil {
			return err
		}
	}

	return nil
}

func (d DefaultSetup) downloadFormulaBundle(url, destPath, zipName, repoName string) (string, error) {
	prompt.Info("Downloading formula...")
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", ErrCreateReqBundle
	}

	if d.edition == api.Team {
		s, err := d.sessionManager.Current()
		if err != nil {
			return "", err
		}
		req.Header.Set(headers.XOrg, s.Organization)
		req.Header.Set(headers.XRepoName, repoName)
		req.Header.Set(headers.Authorization, s.AccessToken)
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusNotFound:
		return "", ErrFormulaBinNotFound
	default:
		return "", ErrUnknownFormulaDownload
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

	prompt.Success("Formula download completed!")
	return file, nil
}

func (d DefaultSetup) downloadConfig(url, destPath, configName, repoName string) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return ErrCreateReqConfig
	}

	if d.edition == api.Team {
		s, err := d.sessionManager.Current()
		if err != nil {
			return err
		}
		req.Header.Set(headers.XOrg, s.Organization)
		req.Header.Set(headers.XRepoName, repoName)
		req.Header.Set(headers.Authorization, s.AccessToken)
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusNotFound:
		return ErrConfigFileNotFound
	default:
		return ErrUnknownConfigFileDownload
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

	return nil
}

func createWorkDir(ritchieHome, binPath string, def Definition) (string, string, error) {
	u := uuid.New().String()
	tDir, tBDir := def.TmpWorkDirPath(ritchieHome, u)

	if err := fileutil.CreateDirIfNotExists(tBDir, 0755); err != nil {
		return "", "", err
	}

	if err := fileutil.CopyDirectory(binPath, tBDir); err != nil {
		return "", "", err
	}

	return tDir, tBDir, nil
}

func unzipFile(filename, destPath string) error {
	if err := fileutil.CreateDirIfNotExists(destPath, 0655); err != nil {
		return err
	}

	if err := fileutil.Unzip(filename, destPath); err != nil {
		return err
	}

	if err := fileutil.RemoveFile(filename); err != nil {
		return err
	}

	return nil
}
