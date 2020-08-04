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
	"fmt"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

type listCredentialCmd struct {
	credential.Settings
}

func NewListCredentialCmd(ss credential.Settings) *cobra.Command {
	l := &listCredentialCmd{ss}

	cmd := &cobra.Command{
		Use:     "credential",
		Short:   "List all credential names and fields.",
		Example: "rit list credential",
		RunE:    l.run(),
	}

	return cmd
}

func printCredentialsTable(fields credential.ListCredDatas) {
	table := uitable.New()
	table.Wrap = true
	table.AddRow(
		prompt.Bold("PROVIDER"),
		prompt.Bold("CONTEXT"),
		prompt.Bold("CREDENTIAL"),
	)

	for _, c := range fields {
		table.AddRow(c.Provider, c.Context, c.Credential)
	}
	if len(table.Rows) < 2 {
		setCmd := prompt.Cyan("rit set credential")
		fmt.Printf("You dont have any credential, use %s\n", setCmd)
	} else {
		fmt.Println(table)
	}
}

func (l listCredentialCmd) run() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		data, err := l.Settings.ReadCredentialsValue(l.CredentialsPath())
		if err != nil {
			return err
		}
		printCredentialsTable(data)
		return nil
	}
}
