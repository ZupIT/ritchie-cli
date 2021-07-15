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

package streams

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// Unzip wrapper for archive.zip
func Unzip(src string, dest string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.Reader.File {

		zippedFile, err := file.Open()
		if err != nil {
			return err
		}

		extractedFilePath := filepath.Join(
			dest,
			file.Name,
		)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(extractedFilePath, file.Mode()); err != nil {
				return err
			}
		} else {
			outputFile, err := os.OpenFile(
				extractedFilePath,
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
				file.Mode(),
			)
			if err != nil {
				return err
			}

			_, err = io.Copy(outputFile, zippedFile)
			if err != nil {
				return err
			}
			outputFile.Close()
		}
		zippedFile.Close()
	}
	return nil
}
