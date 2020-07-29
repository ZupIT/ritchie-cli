package envcredential

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

type CredentialResolver struct {
	credential.Finder
}

const errKeyNotFoundTemplate = `Provider %s has not credencial:%s to fix this verify, config.json of formula`

// NewResolver creates a credential resolver instance of Resolver interface
func NewResolver(cf credential.Finder) CredentialResolver {
	return CredentialResolver{cf}
}

func (c CredentialResolver) Resolve(name string) (string, error) {
	s := strings.Split(name, "_")
	service := strings.ToLower(s[1])
	cred, err := c.Find(service)
	if err != nil {
		return "", err
	}

	k := strings.ToLower(s[2])
	credValue, exist := cred.Credential[k]
	if !exist {
		errMsg := fmt.Sprintf(errKeyNotFoundTemplate, service, strings.ToUpper(name))
		return "", errors.New(prompt.Red(errMsg))
	}
	return credValue, nil
}
