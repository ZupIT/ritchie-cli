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

package formula

const (
	GoLang            = "Go"
	JavaLang          = "Java"
	NodeLang          = "Node"
	PhpLang           = "Php"
	PythonLang        = "Python"
	ShellLang         = "Shell"
	NameBin           = "{{bin-name}}"
	NameModule        = "{{nameModule}}"
	NameBinFirstUpper = "{{bin-name-first-upper}}"
)

var Languages = []string{GoLang, JavaLang, NodeLang, PhpLang, PythonLang, ShellLang}

type LangCreator interface {
	Create(srcDir, pkg, pkgDir, dir string) error
}

type Lang struct {
	Creator
	FileFormat   string
	StartFile    string
	Main         string
	Makefile     string
	WindowsBuild string
	Run          string
	Dockerfile   string
	PackageJson  string
	File         string
	Pkg          string
	Compiled     bool
	UpperCase    bool
}
