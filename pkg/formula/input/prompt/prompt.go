/*
 * Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package prompt

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/pflag"

	"github.com/PaesslerAG/jsonpath"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/input"
	"github.com/ZupIT/ritchie-cli/pkg/stream"

	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

const (
	CachePattern         = "%s/.%s.cache"
	DefaultCacheNewLabel = "Type new value?"
	DefaultCacheQty      = 5
	EmptyItems           = "no items were provided. Please insert a list of items for the input %s in the config.json file of your formula"
)

type InputManager struct {
	envResolvers env.Resolvers
	file         stream.FileWriteReadExister
	prompt.InputList
	prompt.InputText
	input.InputTextDefault
	prompt.InputTextValidator
	prompt.InputBool
	prompt.InputPassword
	prompt.InputMultiselect
}

func NewInputManager(
	env env.Resolvers,
	file stream.FileWriteReadExister,
	inList prompt.InputList,
	inText prompt.InputText,
	inTextValidator prompt.InputTextValidator,
	inDefValue input.InputTextDefault,
	inBool prompt.InputBool,
	inPass prompt.InputPassword,
	inMultiselect prompt.InputMultiselect,
) formula.InputRunner {
	return InputManager{
		envResolvers:       env,
		file:               file,
		InputList:          inList,
		InputText:          inText,
		InputTextValidator: inTextValidator,
		InputTextDefault:   inDefValue,
		InputBool:          inBool,
		InputPassword:      inPass,
		InputMultiselect:   inMultiselect,
	}
}

func (in InputManager) Inputs(cmd *exec.Cmd, setup formula.Setup, f *pflag.FlagSet) error {
	config := setup.Config
	defaultFlag := false
	if f != nil {
		defaultFlag, _ = f.GetBool("default")
	}
	for _, i := range config.Inputs {
		items, err := in.loadItems(i, setup.FormulaPath)
		if err != nil {
			return err
		}
		conditionPass, err := input.VerifyConditional(cmd, i)
		if err != nil {
			return err
		}
		if !conditionPass {
			continue
		}

		inputVal, defaultFlagSet := in.defaultFlag(i, defaultFlag)

		if !defaultFlagSet {
			inputVal, err = in.inputTypeToPrompt(items, i)
			if err != nil {
				return err
			}
		}

		if len(inputVal) != 0 {
			in.persistCache(setup.FormulaPath, inputVal, i, items)
			checkForSameEnv(i.Name)
			input.AddEnv(cmd, i.Name, inputVal)
		}
	}
	return nil
}

func (in InputManager) inputTypeToPrompt(items []string, i formula.Input) (string, error) {
	switch i.Type {

	case input.PassType:
		return in.Password(i.Label, i.Tutorial)

	case input.BoolType:
		valBool, err := in.Bool(i.Label, items, i.Tutorial)
		if err != nil {
			return "", err
		}
		return strconv.FormatBool(valBool), nil

	case input.TextType:
		if items != nil {
			return in.loadInputValList(items, i)
		}
		return in.textValidator(i)

	case input.DynamicType:
		dl, err := in.dynamicList(i.RequestInfo)
		if err != nil {
			return "", err
		}
		return in.List(i.Label, dl, i.Tutorial)
	case input.Multiselect:
		if len(items) > 0 {
			sl, err := in.Multiselect(i)
			if err != nil {
				return "", nil
			}
			return strings.Join(sl, ", "), nil
		}
		return "", fmt.Errorf(EmptyItems, i.Name)

	default:
		return input.ResolveIfReserved(in.envResolvers, i)
	}
}

func checkForSameEnv(envKey string) {
	envKey = strings.ToUpper(envKey)
	if _, exist := os.LookupEnv(envKey); exist {
		warnMsg := fmt.Sprintf(
			"The input param %s has the same name of a machine variable."+
				" It will probably result on unexpect behavior", envKey)
		prompt.Warning(warnMsg)
	}
}

func (in InputManager) defaultFlag(input formula.Input, defaultFlag bool) (string, bool) {
	if defaultFlag && input.Default != "" {
		msg := fmt.Sprintf("Added %s by default: %s", input.Name, input.Default)
		prompt.Info(msg)
		return input.Default, true
	}
	return "", false
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

	inputVal, err := in.List(input.Label, items, input.Tutorial)
	if inputVal == newLabel {
		return in.textValidator(input)
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

func (in InputManager) textValidator(i formula.Input) (string, error) {
	required := input.IsRequired(i)
	var inputVal string
	var err error

	if input.HasRegex(i) {
		inputVal, err = in.textRegexValidator(i, required)
	} else {
		inputVal, err = in.InputTextDefault.Text(i)
	}

	return inputVal, err
}

func (in InputManager) textRegexValidator(input formula.Input, required bool) (string, error) {
	return in.InputTextValidator.Text(input.Label, func(text interface{}) error {
		re := regexp.MustCompile(input.Pattern.Regex)
		if re.MatchString(text.(string)) || (!required && text.(string) == "") {
			return nil
		}

		return errors.New(input.Pattern.MismatchText)
	})
}

func (in InputManager) dynamicList(info formula.RequestInfo) ([]string, error) {
	body, err := makeRequest(info)
	if err != nil {
		return nil, err
	}

	list, err := findValues(info.JsonPath, body)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func makeRequest(info formula.RequestInfo) (interface{}, error) {
	response, err := http.Get(info.Url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, fmt.Errorf("dynamic list request got http status %d expecting some 2xx range", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	requestData := interface{}(nil)

	if err := json.Unmarshal(body, &requestData); err != nil {
		return nil, err
	}

	return requestData, nil
}

func findValues(jsonPath string, requestData interface{}) ([]string, error) {
	dynamicOptions, err := jsonpath.Get(jsonPath, requestData)
	if err != nil {
		return nil, err
	}
	dynamicOptionsStr := fmt.Sprintf("%v", dynamicOptions)
	dynamicOptionsArr := strings.Split(dynamicOptionsStr, " ")

	return dynamicOptionsArr, nil
}
