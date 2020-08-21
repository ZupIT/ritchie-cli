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
	prompt.InputPassword
}

// NewResolver creates a credential resolver instance of Resolver interface
func NewResolver(cf credential.Finder, cs credential.Setter, passwordInput prompt.InputPassword) CredentialResolver {
	return CredentialResolver{cf, cs, passwordInput}
}

func (c CredentialResolver) Resolve(name string) (string, error) {
	s := strings.Split(strings.ToLower(name), "_")
	provider := s[1]
	key := s[2]
	cred, err := c.Find(provider)
	if err != nil {
		// Provider was never set
		cred.Service = provider
		return c.PromptCredential(provider, key, cred)
	}
	credValue, exists := cred.Credential[key]
	if !exists {
		// Provider exists but the expected key doesn't
		return c.PromptCredential(provider, key, cred)
	}

	// Provider and key exist
	return credValue, nil
}

func (c CredentialResolver) PromptCredential(provider, key string, credentialDetail credential.Detail) (string, error) {
	message := fmt.Sprintf("Provider key not found, please provide a value for %s %s: ", provider, key)
	inputVal, err := c.Password(message)
	if err != nil {
		return "", err
	}

	credentialDetail.Credential[key] = inputVal

	if err := c.Set(credentialDetail); err != nil {
		return "", err
	}

	return inputVal, nil
}
