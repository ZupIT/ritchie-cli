package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

type addProviderCmd struct {
	prompt.InputBool
	prompt.InputText
	prompt.InputList
}

func NewAddProviderCmd(
	ib prompt.InputBool,
	it prompt.InputText,
	il prompt.InputList) *cobra.Command {
	a := &addProviderCmd{
		ib,
		it,
		il,
	}

	cmd := &cobra.Command{
		Use:     "provider",
		Short:   "Add new provider",
		Example: "rit add provider",
		RunE:    RunFuncE(a.runStdin(), a.runPrompt()),
	}

	cmd.LocalFlags()
	return cmd
}

func (a addProviderCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {

		addMoreCredentials := true
		provider, _ := a.Text("Enter your provider:", true)

		for addMoreCredentials {

			label, _ := a.Text("Credential key/tag:", true)

			typeList := []string{"text", "password"}
			credentialType, _ := a.List("Want to input the credential as a:", typeList)
			fmt.Println(credentialType, label, provider)

			// c.Inputs = append(c.Inputs, newInput)
			addMoreCredentials, _ = a.Bool("Add one more?", []string{"no", "yes"})
		}

		// homeDir, _ := os.UserHomeDir()
		// providerDir := fmt.Sprintf("%s/.rit/repo/providers.json", homeDir)

		// credentialData, _ := json.Marshal()
		// _ = fileutil.WriteFile(providerDir, credentialData)

		return nil
	}
}

func (a addProviderCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		provider, _ := a.Text("Enter your provider:", true)
		fmt.Println(provider)
		return nil
	}
}
