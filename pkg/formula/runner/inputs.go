package runner

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/formula"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

var ErrInputNotRecognized = errors.New("terminal input not recognized")

type InputManager struct {
	envResolvers env.Resolvers
	prompt.InputList
	prompt.InputText
	prompt.InputBool
	prompt.InputPassword
}

func NewInputManager(
	env env.Resolvers,
	inList prompt.InputList,
	inText prompt.InputText,
	inBool prompt.InputBool,
	inPass prompt.InputPassword) InputManager {
	return InputManager{
		envResolvers:  env,
		InputList:     inList,
		InputText:     inText,
		InputBool:     inBool,
		InputPassword: inPass,
	}
}

func (d InputManager) Inputs(cmd *exec.Cmd, setup formula.Setup, inputType api.TermInputType) error {
	switch inputType {
	case api.Prompt:
		if err := d.fromPrompt(cmd, setup); err != nil {
			return err
		}
	case api.Stdin:
		if err := d.fromStdin(cmd, setup); err != nil {
			return err
		}
	default:
		return ErrInputNotRecognized
	}

	return nil
}

func (d InputManager) fromStdin(cmd *exec.Cmd, setup formula.Setup) error {
	data := make(map[string]interface{})
	if err := stdin.ReadJson(cmd.Stdin, &data); err != nil {
		fmt.Println("The stdin inputs weren't informed correctly. Check the JSON used to execute the command.")
		return err
	}

	config := setup.Config

	for _, input := range config.Inputs {
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
			addEnv(cmd, input.Name, inputVal)
		}
	}
	if len(config.Command) != 0 {
		command := fmt.Sprintf(formula.EnvPattern, formula.CommandEnv, config.Command)
		cmd.Env = append(cmd.Env, command)
	}
	return nil
}

func (d InputManager) fromPrompt(cmd *exec.Cmd, setup formula.Setup) error {
	config := setup.Config
	for _, input := range config.Inputs {
		var inputVal string
		var valBool bool
		items, err := loadItems(input, setup.FormulaPath)
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
		case "password":
			inputVal, err = d.Password(input.Label)
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
			persistCache(setup.FormulaPath, inputVal, input, items)
			addEnv(cmd, input.Name, inputVal)
		}
	}
	if len(config.Command) != 0 {
		command := fmt.Sprintf(formula.EnvPattern, formula.CommandEnv, config.Command)
		cmd.Env = append(cmd.Env, command)
	}
	return nil
}

// addEnv Add environment variable to run formulas.
// add the variable inName=inValue to cmd.Env
func addEnv(cmd *exec.Cmd, inName, inValue string) {
	e := fmt.Sprintf(formula.EnvPattern, strings.ToUpper(inName), inValue)
	cmd.Env = append(cmd.Env, e)
}

func persistCache(formulaPath, inputVal string, input formula.Input, items []string) {
	cachePath := fmt.Sprintf(formula.CachePattern, formulaPath, strings.ToUpper(input.Name))
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
		qtd := formula.DefaultCacheQty
		if input.Cache.Qty != 0 {
			qtd = input.Cache.Qty
		}
		if len(items) > qtd {
			items = items[0:qtd]
		}
		itemsBytes, _ := json.Marshal(items)
		err := fileutil.WriteFile(cachePath, itemsBytes)
		if err != nil {
			fmt.Sprintln("Write file error")
			return
		}

	}
}

func (d InputManager) loadInputValList(items []string, input formula.Input) (string, error) {
	newLabel := formula.DefaultCacheNewLabel
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

func loadItems(input formula.Input, formulaPath string) ([]string, error) {
	if input.Cache.Active {
		cachePath := fmt.Sprintf(formula.CachePattern, formulaPath, strings.ToUpper(input.Name))
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

func (d InputManager) resolveIfReserved(input formula.Input) (string, error) {
	s := strings.Split(input.Type, "_")
	resolver := d.envResolvers[s[0]]
	if resolver != nil {
		return resolver.Resolve(input.Type)
	}
	return "", nil
}
