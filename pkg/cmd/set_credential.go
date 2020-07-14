package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/credential/credsingle"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

var inputTypes = []string{"plain text", "secret"}

// setCredentialCmd type for set credential command
type setCredentialCmd struct {
	credential.Setter
	credential.Settings
	credential.SingleSettings
	prompt.InputText
	prompt.InputBool
	prompt.InputList
	prompt.InputPassword
	prompt.InputMultiline
}

// NewSingleSetCredentialCmd creates a new cmd instance
func NewSingleSetCredentialCmd(
	st credential.Setter,
	ss credential.SingleSettings,
	it prompt.InputText,
	ib prompt.InputBool,
	il prompt.InputList,
	ip prompt.InputPassword) *cobra.Command {
	s := &setCredentialCmd{
		st,
		nil,
		ss,
		it,
		ib,
		il,
		ip,
		nil}
	return newCmd(s)
}

func newCmd(s *setCredentialCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "credential",
		Short: "Set credential",
		Long:  `Set credentials for Github, Gitlab, AWS, UserPass, etc.`,
		RunE:  RunFuncE(s.runStdin(), s.runPrompt()),
	}

	cmd.LocalFlags()

	return cmd
}

func (s setCredentialCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		cred, err := s.promptResolver()
		if err != nil {
			return err
		}

		if err := s.Set(cred); err != nil {
			return err
		}

		prompt.Success(fmt.Sprintf("✔ %s credential saved!", strings.Title(cred.Service)))
		return nil
	}
}

func (s setCredentialCmd) promptResolver() (credential.Detail, error) {
	return s.singlePrompt()
}

func (s setCredentialCmd) singlePrompt() (credential.Detail, error) {
	err := s.WriteDefaultCredentials(credsingle.ProviderPath())
	if err != nil {
		return credential.Detail{}, err
	}

	var credDetail credential.Detail
	cred := credential.Credential{}

	credentials, err := s.ReadCredentials(credsingle.ProviderPath())
	if err != nil {
		return credential.Detail{}, err
	}

	providerArr := credsingle.NewProviderArr(credentials)
	providerChoose, err := s.List("Select your provider", providerArr)
	if err != nil {
		return credDetail, err
	}

	if providerChoose == credsingle.AddNew {
		newProvider, err := s.Text("Define your provider name:", true)
		if err != nil {
			return credDetail, err
		}
		providerArr = append(providerArr, newProvider)

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
			addMoreCredentials, err = s.Bool("Add more credentials to this provider?", []string{"no", "yes"})
			if err != nil {
				return credDetail, err
			}
		}
		credentials[newProvider] = newFields
		err = s.WriteCredentials(credentials, credsingle.ProviderPath())
		if err != nil {
			return credDetail, err
		}

		providerChoose = newProvider
	}

	inputs := credentials[providerChoose]

	for _, i := range inputs {
		var value string
		if i.Type == inputTypes[1] {
			value, err = s.Password(i.Name + ":")
			if err != nil {
				return credDetail, err
			}
		} else {
			value, err = s.Text(i.Name, true)
			if err != nil {
				return credDetail, err
			}
		}
		cred[i.Name] = value
	}
	credDetail.Service = providerChoose
	credDetail.Credential = cred

	return credDetail, nil
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

		prompt.Success(fmt.Sprintf("✔ %s credential saved!", strings.Title(cred.Service)))
		return nil
	}
}

func (s setCredentialCmd) stdinResolver() (credential.Detail, error) {
	var credDetail credential.Detail

	err := stdin.ReadJson(os.Stdin, &credDetail)
	if err != nil {
		prompt.Error(stdin.MsgInvalidInput)
		return credDetail, err
	}
	return credDetail, nil
}

func (s setCredentialCmd) profile(credDetail *credential.Detail) error {
	profiles := map[string]credential.Type{
		"ME (for you)":               credential.Me,
		"OTHER (for another user)":   credential.Other,
		"ORG (for the organization)": credential.Org,
	}
	var types []string
	for k := range profiles {
		types = append(types, k)
	}

	typ, err := s.List("Profile to add credential: ", types)
	if err != nil {
		return err
	}

	if profiles[typ] == credential.Other {
		credDetail.Username, err = s.Text("Username: ", true)
		if err != nil {
			return err
		}
	}

	credDetail.Type = profiles[typ]
	return nil
}
