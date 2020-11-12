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

package template

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const (
	languageDir = "languages"
	rootDir     = "root"
)

var (
	templatePath = []string{"repos", "commons", "templates", "create_formula"}
	errMsg       = `To create a new formula, the commons repository must contain the following structure: 
%s
 └── languages
 └── root

See example: [https://github.com/ZupIT/ritchie-formulas/blob/master/templates/create_formula/README.md]`
)

type Manager interface {
	Languages() ([]string, error)
	LangTemplateFiles(lang string) ([]File, error)
	ResolverNewPath(oldPath, newDir, lang, workspacePath string) (string, error)
	Validate() error
}

type File struct {
	Path  string
	IsDir bool
}

func NewManager(ritchieHome string, dir stream.DirChecker) Manager {
	return DefaultManager{ritchieHome, dir}
}

type DefaultManager struct {
	ritchieHome string
	dir         stream.DirChecker
}

func (tm DefaultManager) templateDir() string {
	return filepath.Join(tm.ritchieHome, filepath.Join(templatePath...))
}

func (tm DefaultManager) Languages() ([]string, error) {
	tplD := filepath.Join(tm.templateDir(), languageDir)

	dirs, err := ioutil.ReadDir(tplD)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, d := range dirs {
		if d.IsDir() {
			result = append(result, d.Name())
		}
	}

	return result, nil
}

func (tm DefaultManager) LangTemplateFiles(lang string) ([]File, error) {
	tplD := tm.templateDir()

	langDir := filepath.Join(tplD, languageDir, lang)

	languageTpl, err := readDirRecursive(langDir)
	if err != nil {
		return nil, err
	}

	rootTplDir := filepath.Join(tplD, rootDir)
	rootTpl, err := readDirRecursive(rootTplDir)
	if err != nil {
		return nil, err
	}

	return append(languageTpl, rootTpl...), nil
}

func readDirRecursive(dir string) ([]File, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var fileNames []File //nolint:prealloc
	for _, f := range files {
		if f.IsDir() {
			dirFiles, err := readDirRecursive(filepath.Join(dir, f.Name()))
			if err != nil {
				return nil, err
			}
			fileNames = append(fileNames, dirFiles...)
		}
		fileNames = append(fileNames, File{
			Path:  filepath.Join(dir, f.Name()),
			IsDir: f.IsDir(),
		})

	}
	return fileNames, nil
}

func (tm DefaultManager) ResolverNewPath(oldPath, formulaPath, lang, workspacePath string) (string, error) {
	tplD := tm.templateDir()
	langTplPath := filepath.Join(tplD, languageDir, lang)
	rootTplPath := filepath.Join(tplD, rootDir)

	if strings.Contains(oldPath, rootTplPath) {
		return strings.Replace(oldPath, rootTplPath, workspacePath, 1), nil
	}

	if strings.Contains(oldPath, langTplPath) {
		return strings.Replace(oldPath, langTplPath, formulaPath, 1), nil
	}

	return "", fmt.Errorf("fail to resolve new Path %s", oldPath)
}

func (tm DefaultManager) Validate() error {
	tplDirPath := tm.templateDir()
	tplLangPath := filepath.Join(tplDirPath, languageDir)
	tplRootPath := filepath.Join(tplDirPath, rootDir)
	invalidErr := fmt.Errorf(errMsg, tplDirPath)
	if !tm.dir.Exists(tplDirPath) || !tm.dir.Exists(tplLangPath) || !tm.dir.Exists(tplRootPath) {
		return invalidErr
	}
	return nil
}
