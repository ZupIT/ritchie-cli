package input

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

type InputTextDefault interface {
	Text(input formula.Input) (string, error)
}

const (
	TextType             = "text"
	ListType             = "list"
	BoolType             = "bool"
	PassType             = "password"
	PathType             = "path"
	DynamicType          = "dynamic"
	MultiselectType      = "multiselect"
	MultiselectSeparator = "|"
	TypeSuffix           = "_type"
)

// addEnv Add environment variable to run formulas.
// add the variable inName=inValue to cmd.Env
func AddEnv(cmd *exec.Cmd, inName, inValue string) {
	e := fmt.Sprintf(formula.EnvPattern, strings.ToUpper(inName), inValue)
	cmd.Env = append(cmd.Env, e)
}

func IsRequired(input formula.Input) bool {
	if input.Required == nil {
		return input.Default == ""
	}

	return *input.Required
}

func HasRegex(input formula.Input) bool {
	return len(input.Pattern.Regex) > 0
}

func containsArray(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func inputConditionVariableExistsOnInputList(variable string, inputList formula.Inputs) bool {
	for _, inputListElement := range inputList {
		if inputListElement.Name == variable {
			return true
		}
	}
	return false
}

func containsSubstring(s string, substr string) bool {
	return strings.Contains(s, substr)
}

func valueContainsAny(inputType string, value string, input string) bool {
	splitInput := strings.Split(input, "|")
	if inputType == MultiselectType {
		splitValue := strings.Split(value, "|")
		for _, i := range splitInput {
			if containsArray(splitValue, i) {
				return true
			}
		}
	} else {
		for _, i := range splitInput {
			if containsSubstring(value, i) {
				return true
			}
		}
	}
	return false
}

func valueContainsAll(inputType string, value string, input string) bool {
	splitInput := strings.Split(input, "|")
	if inputType == MultiselectType {
		splitValue := strings.Split(value, "|")
		for _, v := range splitInput {
			if !containsArray(splitValue, v) {
				return false
			}
		}
	} else {
		for _, v := range splitInput {
			if !containsSubstring(value, v) {
				return false
			}
		}
	}
	return true
}

func valueContainsOnly(inputType string, value string, input string) bool {
	if inputType == MultiselectType {
		splitInput := strings.Split(input, "|")
		splitValue := strings.Split(value, "|")
		if len(splitValue) != len(splitInput) {
			return false
		}
		for _, v := range splitInput {
			if !containsArray(splitValue, v) {
				return false
			}
		}
	} else {
		if !(strings.EqualFold(strings.ToLower(value), strings.ToLower(input))) {
			return false
		}
	}
	return true
}

func VerifyConditional(cmd *exec.Cmd, input formula.Input, inputList formula.Inputs) (bool, error) {

	if input.Condition.Variable == "" {
		return true, nil
	}

	var typeValue string
	var value string

	variable := input.Condition.Variable

	if !inputConditionVariableExistsOnInputList(variable, inputList) {
		return false, fmt.Errorf("config.json: conditional variable %s not found", variable)
	}

	for _, envVal := range cmd.Env {
		components := strings.Split(envVal, "=")
		if strings.ToLower(components[0]) == variable {
			value = components[1]
		} else if strings.ToLower(components[0]) == (variable + TypeSuffix) {
			typeValue = components[1]
		}
	}

	if value == "" {
		return false, nil
	}

	if typeValue == "" {
		return false, fmt.Errorf("config.json: conditional variable %s has no type", variable)
	}

	// Currently using case implementation to avoid adding a dependency module or exposing
	// the code to the risks of running an eval function on a user-defined variable
	// optimizations are welcome, being mindful of the points above
	switch input.Condition.Operator {
	case "==":
		return value == input.Condition.Value, nil
	case "!=":
		return value != input.Condition.Value, nil
	case ">":
		return value > input.Condition.Value, nil
	case ">=":
		return value >= input.Condition.Value, nil
	case "<":
		return value < input.Condition.Value, nil
	case "<=":
		return value <= input.Condition.Value, nil
	case "containsAny":
		return valueContainsAny(typeValue, value, input.Condition.Value), nil
	case "containsAll":
		return valueContainsAll(typeValue, value, input.Condition.Value), nil
	case "containsOnly":
		return valueContainsOnly(typeValue, value, input.Condition.Value), nil
	case "notContainsAny":
		return !valueContainsAny(typeValue, value, input.Condition.Value), nil
	case "notContainsAll":
		return !valueContainsAll(typeValue, value, input.Condition.Value), nil
	default:
		return false, fmt.Errorf(
			"config.json: conditional operator %s not valid. Use any of (==, !=, >, >=, <, <=, containsAny, containsAll, containsOnly, notContainsAny, notContainsAll)",
			input.Condition.Operator,
		)
	}
}
