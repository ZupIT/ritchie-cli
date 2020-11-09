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

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ZupIT/ritchie-cli/pkg/commands"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

func main() {
	startTime := time.Now()
	rootCmd := commands.Build()
	err := rootCmd.Execute()
	if err != nil {
		commands.SendMetric(commands.ExecutionTime(startTime), err.Error())
		errFmt := fmt.Sprintf("%+v", err)
		errFmt = prompt.Red(errFmt)
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", errFmt)
		os.Exit(1)
	}
	commands.SendMetric(commands.ExecutionTime(startTime))
}