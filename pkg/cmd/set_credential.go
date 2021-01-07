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

package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

var inputTypes = []string{"plain text", "secret"}
var inputWay = []string{"type", "file"}

// setCredentialCmd type for set credential command.
type setCredentialCmd struct {
	credential.Setter
	credential.ReaderWriterPather
	prompt.InputText
	prompt.InputBool
	prompt.InputList
	prompt.InputPassword
}

var setCredentialFlags = flags{
	{
		name:        "provider",
		kind:        reflect.String,
		defValue:    "",
		description: "provider name (i.e.: github)",
	},
	{
		name:        "fields",
		kind:        reflect.Slice,
		defValue:    "",
		description: "comma separated list of field names",
	},
	{
		name:        "values",
		kind:        reflect.Slice,
		defValue:    "",
		description: "comma separated list of field values",
	},
}

// NewSetCredentialCmd creates a new cmd instance.
func NewSetCredentialCmd(
	credSetter credential.Setter,
	credFile credential.ReaderWriterPather,
	inText prompt.InputText,
	inBool prompt.InputBool,
	inList prompt.InputList,
	inPass prompt.InputPassword,
) *cobra.Command {
	s := &setCredentialCmd{
		Setter:             credSetter,
		ReaderWriterPather: credFile,
		InputText:          inText,
		InputBool:          inBool,
		InputList:          inList,
		InputPassword:      inPass,
	}

	cmd := &cobra.Command{
		Use:       "credential",
		Short:     "Set credential",
		Long:      `Set credentials for Github, Gitlab, AWS, UserPass, etc.`,
		RunE:      RunFuncE(s.runStdin(), s.runFormula()),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
	}

	addReservedFlags(cmd.Flags(), setCredentialFlags)

	return cmd
}

func (s setCredentialCmd) runFormula() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		cred, err := s.resolveInput(cmd)
		if err != nil {
			return err
		}

		if err := s.Set(cred); err != nil {
			return err
		}
		prompt.Success(fmt.Sprintf("%s credential saved!", strings.Title(cred.Service)))
		prompt.Info("Check your credentials using rit list credential")

		return nil
	}
}

func (s *setCredentialCmd) resolveInput(cmd *cobra.Command) (credential.Detail, error) {
	if IsFlagInput(cmd) {
		return s.resolveFlags(cmd)
	}
	return s.resolvePrompt()
}

func (s *setCredentialCmd) resolvePrompt() (credential.Detail, error) {
	if err := s.WriteDefaultCredentialsFields(s.ProviderPath()); err != nil {
		return credential.Detail{}, err
	}

	var credDetail credential.Detail
	cred := credential.Credential{}

	credentials, err := s.ReadCredentialsFields(s.ProviderPath())
	if err != nil {
		return credential.Detail{}, err
	}

	providerArr := credential.NewProviderArr(credentials)
	providerChoose, err := s.List("Select your provider", providerArr)
	if err != nil {
		return credDetail, err
	}

	if providerChoose == credential.AddNew {
		newProvider, err := s.Text("Define your provider name:", true)
		if err != nil {
			return credDetail, err
		}

		var newFields []credential.Field
		var newField credential.Field
		addMoreCredentials := true
		for addMoreCredentials {
			newField.Name, err = s.Text("Define your field name: (ex.:token, secretAccessKey)", true)
			if err != nil {
				return credDetail, err
			}

			newField.Type, err = s.List("Select your field type:", inputTypes)
			if err != nil {
				return credDetail, err
			}

			newFields = append(newFields, newField)
			addMoreCredentials, err = s.Bool("Add more credentials fields to this provider?", []string{"no", "yes"})
			if err != nil {
				return credDetail, err
			}
		}
		credentials[newProvider] = newFields
		if err = s.WriteCredentialsFields(credentials, s.ProviderPath()); err != nil {
			return credDetail, err
		}

		providerChoose = newProvider
	}

	inputs := credentials[providerChoose]

	inputWayChoose, _ := s.List("Want to enter your credential typing or through a file?", inputWay)
	for _, i := range inputs {
		var value string
		if inputWayChoose == inputWay[1] {
			path, err := s.Text("Enter the credential file path for "+prompt.Cyan(i.Name)+":", true)
			if err != nil {
				return credential.Detail{}, err
			}

			byteValue, err := ioutil.ReadFile(path)
			if err != nil {
				return credential.Detail{}, err
			}
			if len(byteValue) == 0 {
				return credential.Detail{}, prompt.NewError("Empty credential file!")
			}

			cred[i.Name] = string(byteValue)

		} else {
			if i.Type == inputTypes[1] {
				value, err = s.Password(i.Name + ":")
				if err != nil {
					return credDetail, err
				}
			} else {
				value, err = s.Text(i.Name+":", true)
				if err != nil {
					return credDetail, err
				}
			}
			cred[i.Name] = value
		}
	}
	credDetail.Service = providerChoose
	credDetail.Credential = cred

	return credDetail, nil
}

func (s *setCredentialCmd) resolveFlags(cmd *cobra.Command) (credential.Detail, error) {
	provider, err := cmd.Flags().GetString(providerFlagName)
	if err != nil {
		return credential.Detail{}, err
	} else if provider == "" {
		return credential.Detail{}, errors.New("please provide a value for 'provider'")
	}

	fields, err := cmd.Flags().GetStringSlice(fieldsFlagName)
	if err != nil {
		return credential.Detail{}, err
	} else if len(fields) == 0 {
		return credential.Detail{}, errors.New("please provide a value for 'fields'")
	}

	values, err := cmd.Flags().GetStringSlice(valuesFlagName)
	if err != nil {
		return credential.Detail{}, err
	} else if len(values) == 0 {
		return credential.Detail{}, errors.New("please provide a value for 'values'")
	}

	if len(fields) != len(values) {
		return credential.Detail{}, errors.New("number of fields does not match with number of values")
	}

	credentialMap := make(map[string]string)
	for i := 0; i < len(fields); i++ {
		credentialMap[fields[i]] = values[i]
	}

	return credential.Detail{
		Service:    provider,
		Credential: credentialMap,
	}, nil
}

func (s setCredentialCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		cred, err := s.stdinResolver()
		if err != nil {
			return err
		}

		if err := s.Set(cred); err != nil {
			return err
		}

		prompt.Success(fmt.Sprintf("%s credential saved!", strings.Title(cred.Service)))
		prompt.Info("Check your credentials using rit list credential")
		return nil
	}
}

func (s setCredentialCmd) stdinResolver() (credential.Detail, error) {
	var credDetail credential.Detail

	if err := stdin.ReadJson(os.Stdin, &credDetail); err != nil {
		return credDetail, err
	}
	return credDetail, nil
}
