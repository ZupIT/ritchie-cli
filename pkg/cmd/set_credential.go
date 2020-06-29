package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

// setCredentialCmd type for set credential command
type setCredentialCmd struct {
	credential.Setter
	credential.Settings
	edition api.Edition
	prompt.InputText
	prompt.InputBool
	prompt.InputList
	prompt.InputPassword
}

// NewSingleSetCredentialCmd creates a new cmd instance
func NewSingleSetCredentialCmd(
	st credential.Setter,
	it prompt.InputText,
	ib prompt.InputBool,
	il prompt.InputList,
	ip prompt.InputPassword) *cobra.Command {
	s := &setCredentialCmd{st, nil, api.Single, it, ib, il, ip}

	return newCmd(s)
}

// NewTeamSetCredentialCmd creates a new cmd instance
func NewTeamSetCredentialCmd(
	st credential.Setter,
	si credential.Settings,
	it prompt.InputText,
	ib prompt.InputBool,
	il prompt.InputList,
	ip prompt.InputPassword) *cobra.Command {
	s := &setCredentialCmd{st, si, api.Team, it, ib, il, ip}

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

		prompt.Success(fmt.Sprintf("%s credential saved!", strings.Title(cred.Service)))
		return nil
	}
}

func (s setCredentialCmd) promptResolver() (credential.Detail, error) {
	switch s.edition {
	case api.Single:
		return s.singlePrompt()
	case api.Team:
		return s.teamPrompt()
	default:
		return credential.Detail{}, fmt.Errorf(prompt.Red, "invalid CLI build, no edition defined")
	}
}

func (s setCredentialCmd) singlePrompt() (credential.Detail, error) {
	var credDetail credential.Detail
	cred := credential.Credential{}

	credentials, _ := ReadCredentialsJson()

	var credentialList []string
	for _, p := range credentials.SingleCredentials {
		credentialList = append(credentialList, p.Provider)
	}
	providerChoose, _ := s.List("Select your provider", credentialList)

	if providerChoose == "Add new" {
		var c credential.SingleCredential
		newProvider, _ := s.Text("Provider name:", true)
		c.Provider = newProvider
		credentialList = append(credentialList, newProvider)
		addMoreCredentials := true
		for addMoreCredentials {
			var newInput credential.Input
			label, _ := s.Text("Credential key/tag:", true)

			typeList := []string{"text", "password"}
			credentialType, _ := s.List("Want to input the credential as a:", typeList)

			newInput.Type = credentialType
			newInput.Label = label
			c.Inputs = append(c.Inputs, newInput)
			addMoreCredentials, _ = s.Bool("Add one more?", []string{"no", "yes"})

		}
		credentials.SingleCredentials = append(credentials.SingleCredentials, c)


		idk, _ := json.Marshal(credentials)
		home, _ := os.UserHomeDir()
		providerDir := fmt.Sprintf("%s/.rit/repo/providers.json", home)
		_ = ioutil.WriteFile(providerDir, idk, 0644)
		
	}


	providerChoose, _ = s.List("Select your provider", credentialList)

	for _,c := range credentials.SingleCredentials {
		if c.Provider == providerChoose{
			for _, i := range c.Inputs {
				if i.Type == prompt.PasswordType {
					cred[i.Label], _ = s.Password(i.Label)// ra du ken
				} else {
					cred[i.Label], _ = s.Text(i.Label, true)
				}
			}
		}
	}

	credDetail.Service = providerChoose
	credDetail.Credential = cred

	return credDetail, nil
}

func ReadCredentialsJson() (credential.SingleCredentials, error) {
	var credentials credential.SingleCredentials
	home, _ := os.UserHomeDir()
	providerDir := fmt.Sprintf("%s/.rit/repo/providers.json", home)
	jsonFile, err := os.Open(providerDir)
	if err != nil {
		return credentials, err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &credentials)
	if err != nil {
		fmt.Println(err)
	}
	return credentials, err
}


func (s setCredentialCmd) teamPrompt() (credential.Detail, error) {
	var credDetail credential.Detail

	cfg, err := s.Fields()
	if err != nil {
		return credDetail, err
	}

	providers := make([]string, 0, len(cfg))
	for k := range cfg {
		providers = append(providers, k)
	}

	if err := s.profile(&credDetail); err != nil {
		return credDetail, err
	}

	service, err := s.List("Provider: ", providers)
	if err != nil {
		return credDetail, err
	}

	credentials := make(map[string]string)
	fields := cfg[service]
	for _, f := range fields {
		var val string
		var err error
		field := strings.ToLower(f.Name)
		lab := fmt.Sprintf("%s %s: ", strings.Title(service), f.Name)
		if f.Type == prompt.PasswordType {
			val, err = s.Password(lab)
		} else {
			val, err = s.Text(lab, true)
		}
		if err != nil {
			return credDetail, err
		}
		credentials[field] = val
	}

	credDetail.Credential = credentials
	credDetail.Service = service

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

		prompt.Success(fmt.Sprintf("%s credential saved!", strings.Title(cred.Service)))
		return nil
	}
}

func (s setCredentialCmd) stdinResolver() (credential.Detail, error) {
	var credDetail credential.Detail

	if s.edition == api.Single || s.edition == api.Team {

		err := stdin.ReadJson(os.Stdin, &credDetail)
		if err != nil {
			prompt.Error(stdin.MsgInvalidInput)
			return credDetail, err
		}

		return credDetail, nil
	}

	return credDetail, fmt.Errorf(prompt.Red, "invalid CLI build, no edition defined")
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


