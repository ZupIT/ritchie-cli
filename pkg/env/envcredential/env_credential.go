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

package envcredential

import (
	"fmt"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

type CredentialResolver struct {
	credential.Finder
	credential.Setter
}

// NewResolver creates a credential resolver instance of Resolver interface
func NewResolver(cf credential.Finder, cs credential.Setter) CredentialResolver {
	return CredentialResolver{cf, cs}
}

func (c CredentialResolver) Resolve(name string, passwordInput prompt.InputPassword) (string, error) {
	s := strings.Split(strings.ToLower(name), "_")
	service := s[1]
	key := s[2]
	cred, err := c.Find(service)
	if err != nil {
		// Provider was never set
		cred.Service = service
		return c.PromptCredential(service, key, cred, passwordInput)
	}
	credValue, exists := cred.Credential[key]
	if !exists {
		// Provider exists but the expected key doesn't
		return c.PromptCredential(service, key, cred, passwordInput)
	}

	// Provider and key exist
	return credValue, nil
}

func (c CredentialResolver) PromptCredential(
	provider string,
	key string,
	credentialDetail credential.Detail,
	passwordInput prompt.InputPassword,
) (string, error) {
	inputVal, err := passwordInput.Password(
		fmt.Sprintf("Provider key not found, please provide a value for %s %s: ", provider, key),
	)
	if err != nil {
		return "", err
	}

	credentialDetail.Credential[key] = inputVal

	if err := c.Set(credentialDetail); err != nil {
		return "", err
	}

	return inputVal, nil
}
