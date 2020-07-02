package formula

import (
	"os"
	"testing"
)

var def Definition
var home string

func TestMain(m *testing.M) {
	home = os.TempDir()
	def = Definition{
		Path:     "scaffold/coffee-java",
		Bin:      "coffee-java.sh",
		LBin:     "coffee-java.sh",
		MBin:     "coffee-java.sh",
		WBin:     "coffee-java-${so}",
		Bundle:   "commons.zip",
		RepoURL:  "https://localhost:8080/formulas",
		RepoName: "commons",
	}

	os.Exit(m.Run())
}

func TestFormulaPath(t *testing.T) {
	const want = "/tmp/formulas/scaffold/coffee-java"
	got := def.FormulaPath(home)

	if want != got {
		t.Errorf("FormulaPath got %v, want %v", got, want)
	}
}

func TestTmpWorkDirPath(t *testing.T) {
	const hash = "e43c2b35-aa28-4833-b6d3-f1e89691fbd6"
	const wantTmpDir = "/tmp/tmp/e43c2b35-aa28-4833-b6d3-f1e89691fbd6"
	const wantTmpBinDir = "/tmp/tmp/e43c2b35-aa28-4833-b6d3-f1e89691fbd6/scaffold/coffee-java"

	gotTmpDir, gotTmpBinDir := def.TmpWorkDirPath(home, hash)

	if wantTmpDir != gotTmpDir {
		t.Errorf("TmpWorkDirPath got tmp dir %v, want %v", gotTmpDir, wantTmpDir)
	}

	if wantTmpBinDir != gotTmpBinDir {
		t.Errorf("TmpWorkDirPath got tmp bin dir %v, want %v", gotTmpBinDir, wantTmpBinDir)
	}
}

func TestBinPath(t *testing.T) {
	const want = "/tmp/formulas/scaffold/coffee-java/bin"
	formulaPath := def.FormulaPath(home)
	got := def.BinPath(formulaPath)

	if want != got {
		t.Errorf("BinPath got %v, want %v", got, want)
	}
}

func TestBinFilePath(t *testing.T) {
	const want = "/tmp/formulas/scaffold/coffee-java/bin/coffee-java.sh"
	formulaPath := def.FormulaPath(home)
	binPath := def.BinPath(formulaPath)
	got := def.BinFilePath(binPath, def.BinName())

	if want != got {
		t.Errorf("BinFilePath got %v, want %v", got, want)
	}
}

func TestBundleURL(t *testing.T) {
	const want = "https://localhost:8080/formulas/scaffold/coffee-java/commons.zip"
	got := def.BundleURL()
	if want != got {
		t.Errorf("BundleURL got %v, want %v", got, want)
	}
}

func TestConfigPath(t *testing.T) {
	const want = "/tmp/formulas/scaffold/coffee-java/config.json"
	formulaPath := def.FormulaPath(home)
	configName := def.ConfigName()
	got := def.ConfigPath(formulaPath, configName)

	if want != got {
		t.Errorf("ConfigPath got %v, want %v", got, want)
	}
}

func TestConfigURL(t *testing.T) {
	const want = "https://localhost:8080/formulas/scaffold/coffee-java/config.json"
	configName := def.ConfigName()
	got := def.ConfigURL(configName)

	if want != got {
		t.Errorf("ConfigURL got %v, want %v", got, want)
	}
}

func TestConfigName(t *testing.T) {
	tests := []struct {
		name   string
		config string
		want   string
	}{
		{
			name: "default config name",
			want: "config.json",
		},
		{
			name:   "definition config name",
			config: "config-test.json",
			want:   "config-test.json",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			def.Config = test.config
			got := def.ConfigName()

			if test.want != got {
				t.Errorf("ConfigName got %v, want %v", got, test.want)
			}
		})
	}
}

func TestFormulaName(t *testing.T) {
	const want = "create_test"
	create := Create{
		FormulaCmd: "rit create test",
	}

	got := create.FormulaName()

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