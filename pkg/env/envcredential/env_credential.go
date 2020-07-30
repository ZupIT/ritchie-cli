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
	"errors"
	"fmt"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

type CredentialResolver struct {
	credential.Finder
}

const errKeyNotFoundTemplate = `Provider %s has not credential:%s to fix this, verify config.json of formula`

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
