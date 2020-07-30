package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

var inputTypes = []string{"plain text", "secret"}
var inputWay = []string{"file", "type"}

// setCredentialCmd type for set credential command
type setCredentialCmd struct {
	credential.Setter
	credential.ReaderWriterPather
	stream.FileReadExister
	prompt.InputText
	prompt.InputBool
	prompt.InputList
	prompt.InputPassword
}

// NewSetCredentialCmd creates a new cmd instance
func NewSetCredentialCmd(
	credSetter credential.Setter,
	credFile credential.ReaderWriterPather,
	file stream.FileReadExister,
	inText prompt.InputText,
	inBool prompt.InputBool,
	inList prompt.InputList,
	inPass prompt.InputPassword,
) *cobra.Command {
	s := &setCredentialCmd{
		Setter:             credSetter,
		ReaderWriterPather: credFile,
		FileReadExister:    file,
		InputText:          inText,
		InputBool:          inBool,
		InputList:          inList,
		InputPassword:      inPass,
	}

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
		cred, err := s.prompt()
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

func (s setCredentialCmd) prompt() (credential.Detail, error) {
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

	inputWayChoose, _ := s.List("Want to enter your credential through a file or by typing it?", inputWay)
	for _, i := range inputs {
		var value string
		if inputWayChoose == inputWay[1] {
			path, _ := s.Text("Enter the file path for "+i.Name+":", true)
			byteValue, _ := s.FileReadExister.Read(path)
			value = string(byteValue)
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
