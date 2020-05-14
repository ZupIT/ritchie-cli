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
	"os/exec"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

const (
	localTreeFile     = "%s/tree/tree.json"
	nameModule        = "{{nameModule}}"
	nameBin           = "{{bin-name}}"
	nameBinFirstUpper = "{{bin-name-first-upper}}"
)

type DefaultRunner struct {
	ritchieHome  string
	envResolvers env.Resolvers
	client       *http.Client
	treeManager  TreeManager
	prompt.InputList
	prompt.InputText
	prompt.InputBool
}

func NewRunner(
	ritchieHome string,
	er env.Resolvers,
	hc *http.Client,
	tm TreeManager,
	il prompt.InputList,
	it prompt.InputText,
	ib prompt.InputBool) DefaultRunner {
	return DefaultRunner{
		ritchieHome,
		er,
		hc,
		tm,
		il,
		it,
		ib,
	}
}

// Run default implementation of function Manager.Run
func (d DefaultRunner) Run(def Definition, docker bool) error {
	pwd, _ := os.Getwd() // Needs in Post run to copy new files, Must return pwd in Pre run
	formulaPath := def.FormulaPath(d.ritchieHome) // Needs only in Pre run
	config, err := d.loadConfig(formulaPath, def) // Needs in Run inputs, Must return config in pre run?
	if err != nil {
		return err
	}

	binName := def.BinName() // Needs in run, Must return binName in pre run?
	binPath := def.BinPath(formulaPath) // Needs in post run, Must return binPath in pre run?
	binFilePath := def.BinFilePath(binPath, binName) // Needs only in pre run

	if !fileutil.Exists(binFilePath) { // Pre run
		url := def.BundleUrl()
		name := def.BundleName()
		zipFile, err := d.downloadFormulaBundle(url, formulaPath, name)
		if err != nil {
			return err
		}

		if err := d.unzipFile(zipFile, formulaPath); err != nil {
			return err
		}
	} // ---

	tmpDir, tmpBinDir, err := d.createWorkDir(binPath, def) // Pre run, Must return tmpDir to remove in Post run
	if err != nil {
		return err
	}

	defer d.removeWorkDir(tmpDir) // Post run

	if err := os.Chdir(tmpBinDir); err != nil { // Pre run
		return err
	}

	tmpBinFilePath := def.BinFilePath(tmpBinDir, binName) // Needs in run, Must return tmpBinFilePath in Pre run

	cmd := exec.Command(tmpBinFilePath) // Run
	cmd.Stdin = os.Stdin // Run
	cmd.Stdout = os.Stdout // Run
	cmd.Stderr = os.Stderr // Run

	if err := d.inputs(cmd, formulaPath, &config); err != nil { // Run
		return err
	} // ---

	if err := cmd.Start(); err != nil { // Run
		return err
	} // ---

	if err := cmd.Wait(); err != nil { // Run
		return err
	} // ---

	df, err := fileutil.ListNewFiles(binPath, tmpBinDir) // Post run
	if err != nil {
		return err
	} // ---
	if err = fileutil.MoveFiles(tmpBinDir, pwd, df); err != nil { // Post run
		return err
	} // ---
	return nil
}

func (d DefaultRunner) inputs(cmd *exec.Cmd, formulaPath string, config *Config) error {
	for i, input := range config.Inputs {
		var inputVal string
		var valBool bool
		items, err := d.loadItems(input, formulaPath)
		if err != nil {
			return err
		}
		switch iType := input.Type; iType {
		case "text":
			if items != nil {
				inputVal, err = d.loadInputValList(items, input)
			} else {
				validate := input.Default == ""
				inputVal, err = d.Text(input.Label, validate)
				if inputVal == "" {
					inputVal = input.Default
				}
			}
		case "bool":
			valBool, err = d.Bool(input.Label, items)
			inputVal = strconv.FormatBool(valBool)
		default:
			inputVal, err = d.resolveIfReserved(input)
			if err != nil {
				log.Fatalf("Fail to resolve input: %v, verify your credentials. [try using set credential]", input.Type)
			}
		}

		if err != nil {
			return err
		}

		if len(inputVal) != 0 {
			d.persistCache(formulaPath, inputVal, input, items)
			e := fmt.Sprintf(EnvPattern, strings.ToUpper(input.Name), inputVal)
			if i == 0 {
				cmd.Env = append(os.Environ(), e)
			} else {
				cmd.Env = append(cmd.Env, e)
			}
		}
	}
	if len(config.Command) != 0 {
		command := fmt.Sprintf(EnvPattern, CommandEnv, config.Command)
		cmd.Env = append(cmd.Env, command)
	}
	return nil
}

func (d DefaultRunner) loadConfig(formulaPath string, def Definition) (Config, error) {
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

func (d DefaultRunner) removeWorkDir(tmpDir string) {
	if err := fileutil.RemoveDir(tmpDir); err != nil {
		fmt.Sprintln("Error in remove dir")
	}
}

func (d DefaultRunner) persistCache(formulaPath, inputVal string, input Input, items []string) {
	cachePath := fmt.Sprintf(CachePattern, formulaPath, strings.ToUpper(input.Name))
	if input.Cache.Active {
		if items == nil {
			items = []string{inputVal}
		} else {
			for i, item := range items {
				if item == inputVal {
					items = append(items[:i], items[i+1:]...)
					break
				}
			}
			items = append([]string{inputVal}, items...)
		}
		qtd := DefaultCacheQtd
		if input.Cache.Qtd != 0 {
			qtd = input.Cache.Qtd
		}
		if len(items) > qtd {
			items = items[0:qtd]
		}
		itemsBytes, _ := json.Marshal(items)
		err := fileutil.WriteFile(cachePath, itemsBytes)
		if err != nil {
			fmt.Sprintln("Error in WriteFile")
			return
		}

	}
}

func (d DefaultRunner) loadInputValList(items []string, input Input) (string, error) {
	newLabel := DefaultCacheNewLabel
	if input.Cache.Active {
		if input.Cache.NewLabel != "" {
			newLabel = input.Cache.NewLabel
		}
		items = append(items, newLabel)
	}
	inputVal, err := d.List(input.Label, items)
	if inputVal == newLabel {
		validate := len(input.Default) == 0
		inputVal, err = d.Text(input.Label, validate)
		if len(inputVal) == 0 {
			inputVal = input.Default
		}
	}
	return inputVal, err
}

func (d DefaultRunner) loadItems(input Input, formulaPath string) ([]string, error) {
	if input.Cache.Active {
		cachePath := fmt.Sprintf(CachePattern, formulaPath, strings.ToUpper(input.Name))
		if fileutil.Exists(cachePath) {
			fileBytes, err := fileutil.ReadFile(cachePath)
			if err != nil {
				return nil, err
			}
			var items []string
			err = json.Unmarshal(fileBytes, &items)
			if err != nil {
				return nil, err
			}
			return items, nil
		} else {
			itemsBytes, err := json.Marshal(input.Items)
			if err != nil {
				return nil, err
			}
			err = fileutil.WriteFile(cachePath, itemsBytes)
			if err != nil {
				return nil, err
			}
			return input.Items, nil
		}
	} else {
		return input.Items, nil
	}
}

func (d DefaultRunner) resolveIfReserved(input Input) (string, error) {
	s := strings.Split(input.Type, "_")
	resolver := d.envResolvers[s[0]]
	if resolver != nil {
		return resolver.Resolve(input.Type)
	}
	return "", nil
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

	out, err := os.Create(file) // Create an empty .zip file
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil { // Copy body data to empty .zip file
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
