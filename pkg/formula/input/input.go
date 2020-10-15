package input

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

const (
	TextType       = "text"
	BoolType       = "bool"
	PassType       = "password"
	NonInteractive = "non-interactive"
)

// addEnv Add environment variable to run formulas.
// add the variable inName=inValue to cmd.Env
func AddEnv(cmd *exec.Cmd, inName, inValue string) {
	e := fmt.Sprintf(formula.EnvPattern, strings.ToUpper(inName), inValue)
	cmd.Env = append(cmd.Env, e)
}

// Resolve Reserved resolves the environment variables reserved.
// for example: when we add a new credential for inputs that starts with CREDENTIAL
// the function try to find inside the envResolvers for the CREDENTIAL implementation
func ResolveIfReserved(envResolvers env.Resolvers, input formula.Input) (string, error) {
	s := strings.Split(input.Type, "_")
	resolver := envResolvers[s[0]]
	if resolver != nil {
		return resolver.Resolve(input.Type)
	}
	return "", nil
}
