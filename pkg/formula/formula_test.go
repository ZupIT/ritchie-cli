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

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ZupIT/ritchie-cli/pkg/os/osutil"
)

var def Definition
var home string

func TestMain(m *testing.M) {
	home = os.TempDir()
	def = Definition{
		Path:     "scaffold/coffee-java",
		RepoName: "commons",
	}

	os.Exit(m.Run())
}

func TestFormulaPath(t *testing.T) {
	want := filepath.Join(home, "repos", "commons", "scaffold", "coffee-java")
	got := def.FormulaPath(home)

	assert.Equal(t, want, got)
}

func TestTmpWorkDirPath(t *testing.T) {
	want := filepath.Join(home, TmpDir)
	gotTmpDir := def.TmpWorkDirPath(home)

	assert.Contains(t, gotTmpDir, want)
}

func TestBinPath(t *testing.T) {
	want := filepath.Join(home, "repos", "commons", "scaffold", "coffee-java", "bin")
	formulaPath := def.FormulaPath(home)
	got := def.BinPath(formulaPath)

	assert.Equal(t, want, got)
}

func TestBinFilePath(t *testing.T) {
	os := runtime.GOOS
	run := "run.sh"
	if os == osutil.Windows {
		run = "run.bat"
	}
	want := filepath.Join(home, "repos", "commons", "scaffold", "coffee-java", "bin", run)

	formulaPath := def.FormulaPath(home)
	got := def.BinFilePath(formulaPath)

	assert.Equal(t, want, got)
}

func TestFormulaCmdName(t *testing.T) {
	const want = "create test"
	create := Create{
		FormulaCmd: "rit create test",
	}

	got := create.FormulaCmdName()

	assert.Equal(t, want, got)
}

func TestPkgName(t *testing.T) {
	const want = "test"
	create := Create{
		FormulaCmd: "rit create test",
	}

	got := create.PkgName()

	assert.Equal(t, got, want)
}

func TestConfigPath(t *testing.T) {
	want := filepath.Join(home, "repos", "commons", "scaffold", "coffee-java", "config.json")

	got := def.ConfigPath(def.FormulaPath(home))

	assert.Equal(t, got, want)
}
