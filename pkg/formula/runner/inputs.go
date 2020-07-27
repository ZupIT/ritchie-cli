package runner

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

const (
	CachePattern         = "%s/.%s.cache"
	DefaultCacheNewLabel = "Type new value?"
	DefaultCacheQty      = 5
)

var ErrInputNotRecognized = prompt.NewError("terminal input not recognized")

type InputManager struct {
	envResolvers env.Resolvers
	file         stream.FileWriteReadExister
	prompt.InputList
	prompt.InputText
	prompt.InputBool
	prompt.InputPassword
}

func NewInput(
	env env.Resolvers,
	file stream.FileWriteReadExister,
	inList prompt.InputList,
	inText prompt.InputText,
	inBool prompt.InputBool,
	inPass prompt.InputPassword,
) formula.InputRunner {
	return InputManager{
		envResolvers:  env,
		file:          file,
		InputList:     inList,
		InputText:     inText,
		InputBool:     inBool,
		InputPassword: inPass,
	}
}

func (in InputManager) Inputs(cmd *exec.Cmd, setup formula.Setup, inputType api.TermInputType) error {
	switch inputType {
	case api.Prompt:
		if err := in.fromPrompt(cmd, setup); err != nil {
			return err
		}
	case api.Stdin:
		if err := in.fromStdin(cmd, setup); err != nil {
			return err
		}
	default:
		return ErrInputNotRecognized
	}

	return nil
}

func (in InputManager) fromStdin(cmd *exec.Cmd, setup formula.Setup) error {
	data := make(map[string]interface{})
	if err := stdin.ReadJson(cmd.Stdin, &data); err != nil {
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
			inputVal, err = in.resolveIfReserved(input)
			if err != nil {
				return err
			}
		}

		if len(inputVal) != 0 {
			addEnv(cmd, input.Name, inputVal)
		}
	}
	return nil
}

func (in InputManager) fromPrompt(cmd *exec.Cmd, setup formula.Setup) error {
	config := setup.Config
	for _, input := range config.Inputs {
		var inputVal string
		var valBool bool
		items, err := in.loadItems(input, setup.FormulaPath)
		if err != nil {
			return err
		}
		switch iType := input.Type; iType {
		case "text":
			if items != nil {
				inputVal, err = in.loadInputValList(items, input)
			} else {
				validate := input.Default == ""
				inputVal, err = in.Text(input.Label, validate)
				if inputVal == "" {
					inputVal = input.Default
				}
			}
		case "bool":
			valBool, err = in.Bool(input.Label, items)
			inputVal = strconv.FormatBool(valBool)
		case "password":
			inputVal, err = in.Password(input.Label)
		default:
			inputVal, err = in.resolveIfReserved(input)
		}

		if err != nil {
			return err
		}

		if len(inputVal) != 0 {
			in.persistCache(setup.FormulaPath, inputVal, input, items)
			addEnv(cmd, input.Name, inputVal)
		}
	}
	return nil
}

// addEnv Add environment variable to run formulas.
// add the variable inName=inValue to cmd.Env
func addEnv(cmd *exec.Cmd, inName, inValue string) {
	e := fmt.Sprintf(formula.EnvPattern, strings.ToUpper(inName), inValue)
	cmd.Env = append(cmd.Env, e)
}

func (in InputManager) persistCache(formulaPath, inputVal string, input formula.Input, items []string) {
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
		qtd := DefaultCacheQty
		if input.Cache.Qty != 0 {
			qtd = input.Cache.Qty
		}
		if len(items) > qtd {
			items = items[0:qtd]
		}
		itemsBytes, _ := json.Marshal(items)
		if err := in.file.Write(cachePath, itemsBytes); err != nil {
			fmt.Sprintln("Write file error")
			return
		}

	}
}

func (in InputManager) loadInputValList(items []string, input formula.Input) (string, error) {
	newLabel := DefaultCacheNewLabel
	if input.Cache.Active {
		if input.Cache.NewLabel != "" {
			newLabel = input.Cache.NewLabel
		}
		items = append(items, newLabel)
	}
	inputVal, err := in.List(input.Label, items)
	if inputVal == newLabel {
		validate := len(input.Default) == 0
		inputVal, err = in.Text(input.Label, validate)
		if len(inputVal) == 0 {
			inputVal = input.Default
		}
	}
	return inputVal, err
}

func (in InputManager) loadItems(input formula.Input, formulaPath string) ([]string, error) {
	if input.Cache.Active {
		cachePath := fmt.Sprintf(CachePattern, formulaPath, strings.ToUpper(input.Name))
		if in.file.Exists(cachePath) {
			fileBytes, err := in.file.Read(cachePath)
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
			if err = in.file.Write(cachePath, itemsBytes); err != nil {
				return nil, err
			}
			return input.Items, nil
		}
	} else {
		return input.Items, nil
	}
}

func (in InputManager) resolveIfReserved(input formula.Input) (string, error) {
	s := strings.Split(input.Type, "_")
	resolver := in.envResolvers[s[0]]
	if resolver != nil {
		return resolver.Resolve(input.Type)
	}
	return "", nil
}
