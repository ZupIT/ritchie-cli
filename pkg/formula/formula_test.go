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
	"strings"
	"testing"

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

	if want != got {
		t.Errorf("FormulaPath got %v, want %v", got, want)
	}
}

func TestTmpWorkDirPath(t *testing.T) {
	want := filepath.Join(home, TmpDir)
	gotTmpDir := def.TmpWorkDirPath(home)

	if !strings.Contains(gotTmpDir, want) {
		t.Errorf("TmpWorkDirPath got tmp dir %v, want some string that contains %v", gotTmpDir, want)
	}
}

func TestBinPath(t *testing.T) {
	want := filepath.Join(home, "repos", "commons", "scaffold", "coffee-java", "bin")
	formulaPath := def.FormulaPath(home)
	got := def.BinPath(formulaPath)

	if want != got {
		t.Errorf("BinPath got %v, want %v", got, want)
	}
}

func TestBinFilePath(t *testing.T) {
	os := runtime.GOOS
	run := "run.sh"
	if os == osutil.Windows {
		run =  "run.bat"
	}
	want := filepath.Join(home, "repos", "commons", "scaffold", "coffee-java", "bin", run)

	formulaPath := def.FormulaPath(home)
	got := def.BinFilePath(formulaPath)

	if want != got {
		t.Errorf("BinFilePath got %v, want %v", got, want)
	}
}

func TestFormulaCmdName(t *testing.T) {
	const want = "create test"
	create := Create{
		FormulaCmd: "rit create test",
	}

	got := create.FormulaCmdName()

	if want != got {
		t.Errorf("FormulaName got %v, want %v", got, want)
	}
}

func TestPkgName(t *testing.T) {
	const want = "test"
	create := Create{
		FormulaCmd: "rit create test",
	}

	got := create.PkgName()

	if want != got {
		t.Errorf("PkgName got %v, want %v", got, want)
	}
}

func TestConfigPath(t *testing.T) {
	 want := filepath.Join(home, "repos", "commons", "scaffold", "coffee-java", "config.json")

	got := def.ConfigPath(def.FormulaPath(home))

	if want != got {
		t.Errorf("TestConfigPath got %v, want %v", got, want)
	}
}
