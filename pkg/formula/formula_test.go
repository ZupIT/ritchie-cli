package formula

import (
	"os"
	"path"
	"strings"
	"testing"
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
	const want = "/tmp/repos/commons/scaffold/coffee-java"
	got := def.FormulaPath(home)

	if want != got {
		t.Errorf("FormulaPath got %v, want %v", got, want)
	}
}

func TestTmpWorkDirPath(t *testing.T) {

	want := path.Join(home, TmpDir)
	gotTmpDir := def.TmpWorkDirPath(home)

	if !strings.Contains(gotTmpDir, want) {
		t.Errorf("TmpWorkDirPath got tmp dir %v, want some string that contains %v", gotTmpDir, want)
	}
}

func TestBinPath(t *testing.T) {
	const want = "/tmp/repos/commons/scaffold/coffee-java/bin"
	formulaPath := def.FormulaPath(home)
	got := def.BinPath(formulaPath)

	if want != got {
		t.Errorf("BinPath got %v, want %v", got, want)
	}
}

func TestBinFilePath(t *testing.T) {
	const want = "/tmp/repos/commons/scaffold/coffee-java/bin/run.sh"
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

	const want = "/tmp/repos/commons/scaffold/coffee-java/config.json"

	got := def.ConfigPath(def.FormulaPath(home))

	if want != got {
		t.Errorf("TestConfigPath got %v, want %v", got, want)
	}
}
