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
	TextType    = "text"
	ListType    = "list"
	BoolType    = "bool"
	PassType    = "password"
	DynamicType = "dynamic"
	Multiselect = "multiselect"
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

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func valueContainsAny(value string, input string) bool {
	splitValue := strings.Split(value, "|")
	splitInput := strings.Split(input, "|")
	for _, v := range splitInput {
		if contains(splitValue, v) {
			return true
		}
	}
	return false
}

func valueContainsAll(value string, input string) bool {
	splitValue := strings.Split(value, "|")
	splitInput := strings.Split(input, "|")
	for _, v := range splitInput {
		if !contains(splitValue, v) {
			return false
		}
	}
	return true
}

func valueContainsOnly(value string, input string) bool {
	splitValue := strings.Split(value, "|")
	splitInput := strings.Split(input, "|")
	if len(splitValue) != len(splitInput) {
		return false
	}
	for _, v := range splitInput {
		if !contains(splitValue, v) {
			return false
		}
	}
	return true
}

func VerifyConditional(cmd *exec.Cmd, input formula.Input) (bool, error) {
	if input.Condition.Variable == "" {
		return true, nil
	}

	var value string
	variable := input.Condition.Variable
	for _, envVal := range cmd.Env {
		components := strings.Split(envVal, "=")
		if strings.ToLower(components[0]) == variable {
			value = components[1]
			break
		}
	}
	if value == "" {
		return false, fmt.Errorf("config.json: conditional variable %s not found", variable)
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
		return valueContainsAny(value, input.Condition.Value), nil
	case "containsAll":
		return valueContainsAll(value, input.Condition.Value), nil
	case "containsOnly":
		return valueContainsOnly(value, input.Condition.Value), nil
	case "notContainsAny":
		return !valueContainsAny(value, input.Condition.Value), nil
	case "notContainsAll":
		return !valueContainsAll(value, input.Condition.Value), nil
	default:
		return false, fmt.Errorf(
			"config.json: conditional operator %s not valid. Use any of (==, !=, >, >=, <, <=, containsAny, containsAll, containsOnly, notContainsAny, notContainsAll)",
			input.Condition.Operator,
		)
	}
}
