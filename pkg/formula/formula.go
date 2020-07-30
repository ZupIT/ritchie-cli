package formula

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/uuid"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/os/osutil"
)

const (
	ReposDir      = "repos"
	TmpDir        = "tmp"
	DefaultConfig = "config.json"
	PwdEnv        = "PWD"
	CPwdEnv       = "CURRENT_PWD"
	BinUnix       = "run.sh"
	BinWindows    = "run.bat"
	BinDir        = "bin"
	EnvPattern    = "%s=%s"
	TreePath      = "/tree/tree.json"
	MakefilePath  = "/Makefile"
)

type (
	Input struct {
		Name    string   `json:"name"`
		Type    string   `json:"type"`
		Default string   `json:"default"`
		Label   string   `json:"label"`
		Items   []string `json:"items"`
		Cache   Cache    `json:"cache"`
	}

	Cache struct {
		Active   bool   `json:"active"`
		Qty      int    `json:"qty"`
		NewLabel string `json:"newLabel"`
	}
	Create struct {
		FormulaCmd    string `json:"formulaCmd"`
		Lang          string `json:"lang"`
		WorkspacePath string `json:"workspacePath"`
		FormulaPath   string `json:"formulaPath"`
	}

	Config struct {
		DockerIB string  `json:"dockerImageBuilder"`
		Inputs   []Input `json:"inputs"`
	}

	// Definition type that represents a Formula
	Definition struct {
		Path     string
		RepoName string
	}

	Setup struct {
		Pwd         string
		FormulaPath string
		BinName     string
		BinPath     string
		TmpDir      string
		Config      Config
		ContainerId string
	}
	Help struct {
		Short string `json:"short"`
		Long  string `json:"long"`
	}
)

type PreRunner interface {
	PreRun(def Definition, docker bool) (Setup, error)
}

type Runner interface {
	Run(def Definition, inputType api.TermInputType, docker bool) error
}

type PostRunner interface {
	PostRun(p Setup, docker bool) error
}

type InputRunner interface {
	Inputs(cmd *exec.Cmd, setup Setup, inputType api.TermInputType) error
}

type Creator interface {
	Create(cf Create) error
}

type MakeBuilder interface {
	Build(formulaPath string) error
}

type BatBuilder interface {
	Build(formulaPath string) error
}

type DockerBuilder interface {
	Build(formulaPath, dockerImg string) error
}

type LocalBuilder interface {
	Build(workspacePath, formulaPath string) error
}

type Watcher interface {
	Watch(workspacePath, formulaPath string)
}

type CreateBuilder interface {
	Creator
	LocalBuilder
}

// FormulaPath builds the formula path from ritchie home
func (d *Definition) FormulaPath(home string) string {
	return filepath.Join(home, ReposDir, d.RepoName, d.Path)
}

// TmpWorkDirPath builds the tmp paths to run formula, first parameter is tmpDir created
// second parameter is tmpBinDir
func (d *Definition) TmpWorkDirPath(home string) string {
	u := uuid.New().String()
	return filepath.Join(home, TmpDir, u)
}

// BinFilePath builds the bin file path from formula path
func (d *Definition) BinFilePath(fPath string) string {
	return filepath.Join(fPath, BinDir, d.BinName())
}

func (d *Definition) BinName() string {
	bName := BinUnix
	if runtime.GOOS == osutil.Windows {
		bName = BinWindows
	}
	return bName
}

// BinFilePath builds the bin file path from formula path
func (d *Definition) BinPath(fPath string) string {
	return filepath.Join(fPath, BinDir)
}

// ConfigPath builds the config path from formula path and config name
func (d *Definition) ConfigPath(formulaPath string) string {
	return filepath.Join(formulaPath, DefaultConfig)
}

// FormulaName remove rit from formulaCmd
func (c Create) FormulaCmdName() string {
	d := strings.Split(c.FormulaCmd, " ")
	return strings.Join(d[1:], " ")
}

func (c Create) PkgName() string {
	d := strings.Split(c.FormulaCmd, " ")
	return d[len(d)-1]
}
