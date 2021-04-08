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
)

// addEnv Add environment variable to run formulas.
// add the variable inName=inValue to cmd.Env
func AddEnv(cmd *exec.Cmd, inName, inValue string) {
	e := fmt.Sprintf(formula.EnvPattern, strings.ToUpper(inName), inValue)
	cmd.Env = append(cmd.Env, e)
}

func IsRequired(input formula.Input) bool {
	if input.Required == nil {
		return false
	}

	return *input.Required
}

func HasRegex(input formula.Input) bool {
	return len(input.Pattern.Regex) > 0
}

func inputConditionVariableExistsOnInputList(variable string, inputList formula.Inputs) bool {
	for _, inputListElement := range inputList {
		if inputListElement.Name == variable {
			return true
		}
	}
	return false
}

func VerifyConditional(cmd *exec.Cmd, input formula.Input, inputList formula.Inputs) (bool, error) {
	if input.Condition.Variable == "" {
		return true, nil
	}

	variable := input.Condition.Variable

	if !inputConditionVariableExistsOnInputList(variable, inputList) {
		return false, fmt.Errorf("config.json: conditional variable %s not found", variable)
	}

	var value string
	for _, envVal := range cmd.Env {
		components := strings.Split(envVal, "=")
		if strings.ToLower(components[0]) == variable {
			value = components[1]
			break
		}
	}

	if value == "" {
		return false, nil
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
	default:
		return false, fmt.Errorf(
			"config.json: conditional operator %s not valid. Use any of (==, !=, >, >=, <, <=)",
			input.Condition.Operator,
		)
	}
}
