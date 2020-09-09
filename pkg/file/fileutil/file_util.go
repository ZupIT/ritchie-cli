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

package fileutil

import (
	"fmt"
	"io/ioutil"
	"os"
)

// Exists check if file exists
// Deprecated: use the stream package to work with files and directories
func Exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

// CreateDirIfNotExists creates dir if not exists
// Deprecated: use the stream package to work with files and directories
func CreateDirIfNotExists(dir string, perm os.FileMode) error {
	if Exists(dir) {
		return nil
	}

	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}

// RemoveDir removes path and any children it contains.
// Deprecated: use the stream package to work with files and directories
func RemoveDir(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	return nil
}

// ReadFile wrapper for ioutil.ReadFile
// Deprecated: use the stream package to work with files and directories
func ReadFile(path string) ([]byte, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return f, err
}

// WriteFilePerm wrapper for ioutil.WriteFile
// Deprecated: use the stream package to work with files and directories
func WriteFilePerm(path string, content []byte, perm int32) error {
	return ioutil.WriteFile(path, content, os.FileMode(perm))
}

// Read all files and dir in current directory
func readFilesDir(path string) ([]os.FileInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fl, err := f.Readdir(-1)
	f.Close()
	return fl, err
}

// List new files in nPath differing of oPath
func ListNewFiles(oPath, nPath string) ([]string, error) {
	of, err := readFilesDir(oPath)
	if err != nil {
		return nil, err
	}
	nf, err := readFilesDir(nPath)
	if err != nil {
		return nil, err
	}
	control := make(map[string]int)
	for _, file := range of {
		control[file.Name()]++
	}
	var new []string
	for _, file := range nf {
		if control[file.Name()] == 0 {
			new = append(new, file.Name())
		}
	}
	return new, nil
}
