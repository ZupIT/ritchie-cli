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
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func Test_ListCredentialCmd(t *testing.T) {
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	homeDir, _ := os.UserHomeDir()
	credSettings := credential.NewSettings(fileManager, dirManager, homeDir)

	t.Run("Success case", func(t *testing.T) {
		o := NewListCredentialCmd(credSettings)
		if err := o.Execute();err !=nil{
			t.Errorf("Test_ListCredentialCmd error = %s", err)
		}
	})

}
