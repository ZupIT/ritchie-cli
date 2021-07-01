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

package runner

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	OutputFile = "output.json"
)

func Output(data []string, binPath string) error {
	sanitizeData := []string{}
	flattenData := []string{}
	transformData := make(map[string]string)

	for i := range data {
		output := strings.Split(data[i], " ")[1:]
		newOutput := strings.Join(output, " ")
		sanitizeData = append(sanitizeData, newOutput)

		for el := range sanitizeData {
			element := strings.Split(sanitizeData[el], " ")
			for j := range element {
				flattenData = append(flattenData, element[j])
			}
		}
	}

	for i := range flattenData {
		test := strings.Split(flattenData[i], "=")
		key := test[0]
		value := test[1]
		transformData[key] = value
	}

	outputPath := filepath.Join(binPath, OutputFile)

	outputJSON, err := json.MarshalIndent(transformData, "", "\t")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(outputPath, outputJSON, os.ModePerm); err != nil {
		return err
	}

	return nil
}
