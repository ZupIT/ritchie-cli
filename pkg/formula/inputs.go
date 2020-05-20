package formula

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

type InputManager struct {
	envResolvers env.Resolvers
	prompt.InputList
	prompt.InputText
	prompt.InputBool
}

func NewInputManager(
	env env.Resolvers,
	inList prompt.InputList,
	inText prompt.InputText,
	inBool prompt.InputBool) InputManager {
	return InputManager{
		envResolvers: env,
		InputList:    inList,
		InputText:    inText,
		InputBool:    inBool,
	}
}

func (d InputManager) Inputs(cmd *exec.Cmd, setup Setup, inputType api.TermInputType, docker bool) error {
	switch inputType {
	case api.Prompt:
		if err := d.fromPrompt(cmd, setup, docker); err != nil {
			return err
		}
	case api.Stdin:
		if err := d.fromStdin(cmd, setup, docker); err != nil {
			return err
		}
	default:
		return fmt.Errorf("terminal input (%v) not recongnized", inputType)
	}

	return nil
}

func (d InputManager) fromStdin(cmd *exec.Cmd, setup Setup, docker bool) error {
	data := make(map[string]interface{})
	if err := stdin.ReadJson(os.Stdin, &data); err != nil {
		fmt.Println("The stdin inputs weren't informed correctly. Check the JSON used to execute the command.")
		return err
	}

	config := setup.config

	for i, input := range config.Inputs {
		var inputVal string
		var err error
		switch iType := input.Type; iType {
		case "text", "bool":
			inputVal = fmt.Sprintf("%v", data[input.Name])
		default:
			inputVal, err = d.resolveIfReserved(input)
			if err != nil {
				log.Fatalf("Fail to resolve input: %v, verify your credentials. [try using set credential]", input.Type)
				return err
			}
		}

		if len(inputVal) != 0 {
			if err := addEnv(cmd, setup.pwd, input.Name, inputVal, i, docker); err != nil {
				return err
			}
		}
	}
	if len(config.Command) != 0 {
		command := fmt.Sprintf(EnvPattern, CommandEnv, config.Command)
		cmd.Env = append(cmd.Env, command)
	}
	return nil
}

func (d InputManager) fromPrompt(cmd *exec.Cmd, setup Setup, docker bool) error {
	config := setup.config
	for i, input := range config.Inputs {
		var inputVal string
		var valBool bool
		items, err := loadItems(input, setup.formulaPath)
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
			persistCache(setup.formulaPath, inputVal, input, items)
			if err := addEnv(cmd, setup.pwd, input.Name, inputVal, i, docker); err != nil {
				return err
			}
		}
	}
	if len(config.Command) != 0 {
		command := fmt.Sprintf(EnvPattern, CommandEnv, config.Command)
		cmd.Env = append(cmd.Env, command)
	}
	return nil
}

// addEnv Add environment variable to run formulas.
// If docker is true, create a file named .env and add the variable inName=inValue.
// If docker is false, add the variable inName=inValue to cmd
func addEnv(cmd *exec.Cmd, pwd, inName, inValue string, index int, docker bool) error {
	e := fmt.Sprintf(EnvPattern, strings.ToUpper(inName), inValue)
	if docker {
		if !fileutil.Exists(envFile) {
			pwdEnv := fmt.Sprintf(EnvPattern, PwdEnv, pwd) // Add "pwd" to use in formulas that need it
			file := fmt.Sprintf("%s\n%s\n", pwdEnv, e)
			if err := fileutil.WriteFile(envFile, []byte(file)); err != nil {
				return err
			}
		} else {
			if err := fileutil.AppendFileData(envFile, []byte(e+"\n")); err != nil {
				return err
			}
		}
	} else {
		if index == 0 {
			pwdEnv := fmt.Sprintf(EnvPattern, PwdEnv, pwd)
			cmd.Env = append(cmd.Env, pwdEnv) // Add "pwd" to use in formulas that need it
			cmd.Env = append(cmd.Env, e)
		} else {
			cmd.Env = append(cmd.Env, e)
		}
	}

	return nil
}

func persistCache(formulaPath, inputVal string, input Input, items []string) {
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

func (d InputManager) loadInputValList(items []string, input Input) (string, error) {
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

func loadItems(input Input, formulaPath string) ([]string, error) {
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

func (d InputManager) resolveIfReserved(input Input) (string, error) {
	s := strings.Split(input.Type, "_")
	resolver := d.envResolvers[s[0]]
	if resolver != nil {
		return resolver.Resolve(input.Type)
	}
	return "", nil
}
