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
 └── language 1
 	└── template 1
	└── template 2
 └── language 2
 	└── template 1
 └── root

And the template must contain "build.bat" and "build.sh" files 
See example: [https://github.com/ZupIT/ritchie-templates]`
)

type Manager interface {
	Languages() ([]string, error)
	Templates(lang string) ([]string, error)
	TemplateFiles(lang, tpl string) ([]File, error)
	ResolverNewPath(oldPath, newDir, lang, tpl, workspacePath string) (string, error)
	Validate(lang, tpl string) error
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

func (tm DefaultManager) Templates(lang string) ([]string, error) {
	tplD := tm.templateDir()

	langDir := filepath.Join(tplD, languageDir, lang)

	dirs, err := ioutil.ReadDir(langDir)
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

func (tm DefaultManager) TemplateFiles(lang, tpl string) ([]File, error) {
	var tplDir string
	tplD := tm.templateDir()

	if tpl == "src" {
		tplDir = filepath.Join(tplD, languageDir, lang)
	} else {
		tplDir = filepath.Join(tplD, languageDir, lang, tpl)
	}

	template, err := readDirRecursive(tplDir)
	if err != nil {
		return nil, err
	}

	rootTplDir := filepath.Join(tplD, rootDir)
	rootTpl, err := readDirRecursive(rootTplDir)
	if err != nil {
		return nil, err
	}

	return append(template, rootTpl...), nil
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

func (tm DefaultManager) ResolverNewPath(oldPath, formulaPath, lang, tpl, workspacePath string) (string, error) {
	tplD := tm.templateDir()
	tplTypePath := filepath.Join(tplD, languageDir, lang, tpl)
	rootTplPath := filepath.Join(tplD, rootDir)

	if strings.Contains(oldPath, rootTplPath) {
		return strings.Replace(oldPath, rootTplPath, workspacePath, 1), nil
	}

	if strings.Contains(oldPath, tplTypePath) {
		return strings.Replace(oldPath, tplTypePath, formulaPath, 1), nil
	}

	return "", fmt.Errorf("fail to resolve new Path %s", oldPath)
}

func (tm DefaultManager) Validate(lang, tpl string) error {
	var tplDir string
	tplDirPath := tm.templateDir()
	tplLangPath := filepath.Join(tplDirPath, languageDir)
	tplRootPath := filepath.Join(tplDirPath, rootDir)
	invalidErr := fmt.Errorf(errMsg, tplDirPath)

	if tpl == "src" {
		tplDir = filepath.Join(tplDirPath, languageDir, lang)
	} else {
		tplDir = filepath.Join(tplDirPath, languageDir, lang, tpl)
	}

	if !tm.dir.Exists(tplDirPath) || !tm.dir.Exists(tplLangPath) || !tm.dir.Exists(tplRootPath) || !isValidTemplate(tplDir) {
		return invalidErr
	}
	return nil
}

func isValidTemplate(repoPath string) bool {
	hasBuildBat := false
	hasBuildSh := false
	files, err := ioutil.ReadDir(repoPath)
	if err != nil {
		return false
	}

	for _, file := range files {
		if file.Name() == "build.bat" {
			hasBuildBat = true
		}
		if file.Name() == "build.sh" {
			hasBuildSh = true
		}
	}

	if hasBuildBat && hasBuildSh {
		return true
	} else {
		return false
	}

}
