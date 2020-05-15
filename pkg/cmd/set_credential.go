package cmd

import (
	"errors"
	"fmt"
	"log"
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
		RunE: RunFuncE(s.runFunc(), s.runStdin()),
	}

	cmd.LocalFlags()

	return cmd
}

func (s setCredentialCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		cred, err := s.promptResolver()
		if err != nil {
			return err
		}

		if err := s.Set(cred); err != nil {
			return err
		}

		log.Println(fmt.Sprintf("%s credential saved!", strings.Title(cred.Service)))
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
		return credential.Detail{}, errors.New("invalid CLI build, no edition defined")
	}
}

func (s setCredentialCmd) singlePrompt() (credential.Detail, error) {
	var credDetail credential.Detail

	provider, err := s.Text("Provider: ", true)
	if err != nil {
		return credDetail, err
	}

	cred := credential.Credential{}
	addMore := true
	for addMore {
		kv, err := s.Text("Type your credential using the format key=value (e.g. email=example@example.com): ", true)
		if err != nil {
			return credDetail, err
		}

		pair := strings.Split(kv, "=")
		if s := validate(pair); s != "" {
			fmt.Println(s)
			continue
		}

		cred[pair[0]] = pair[1]

		addMore, err = s.Bool("Add more fields?", []string{"yes", "no"})
		if err != nil {
			return credDetail, err
		}
	}

	credDetail.Service = provider
	credDetail.Credential = cred

	return credDetail, nil
}

func validate(pair []string) string {
	if len(pair) < 2 {
		return "Invalid key value credential"
	}

	if strings.TrimSpace(pair[0]) == "" {
		return "The key must not be empty."
	}

	if strings.TrimSpace(pair[1]) == "" {
		return "The value must not be empty."
	}

	return ""
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

	typ, err := s.List("Profile: ", []string{credential.Me, credential.Other})
	if err != nil {
		return credDetail, err
	}

	username := "me"
	if typ == credential.Other {
		username, err = s.Text("Username: ", true)
		if err != nil {
			return credDetail, err
		}
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

	credDetail.Username = username
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

		log.Println(fmt.Sprintf("%s credential saved!", strings.Title(cred.Service)))
		return nil
	}
}

func (s setCredentialCmd) stdinResolver() (credential.Detail, error) {
	var credDetail credential.Detail

	if s.edition == api.Single || s.edition == api.Team {

		err := stdin.ReadJson(os.Stdin, &credDetail)
		if err != nil {
			fmt.Println("The STDIN inputs weren't informed correctly. Check the JSON used to execute the command.")
			return credDetail, err
		}

		return credDetail, nil
	}

	return credDetail, errors.New("invalid CLI build, no edition defined")
}
