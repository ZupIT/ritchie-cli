package formula

import (
	"fmt"
	"runtime"
	"strings"
)

const (
	PathPattern           = "%s/formulas/%s"
	TmpDirPattern         = "%s/tmp/%s"
	TmpBinDirPattern      = "%s/tmp/%s/%s"
	DefaultConfig         = "config.json"
	ConfigPattern         = "%s/%s"
	CommandEnv            = "COMMAND"
	BinPattern            = "%s%s"
	BinPathPattern        = "%s/bin"
	windows               = "windows"
	EnvPattern            = "%s=%s"
	CachePattern          = "%s/.%s.cache"
	DefaultCacheNewLabel  = "Type new value?"
	DefaultCacheQtd       = 5
	FormCreatePathPattern = "%s/ritchie-formulas-local"
	TreeCreatePathPattern = "%s/tree/tree.json"
	Makefile              = "Makefile"
)

// Config type that represents formula config
type Config struct {
	Name        string  `json:"name"`
	Command     string  `json:"command"`
	Description string  `json:"description"`
	Language    string  `json:"language"`
	Inputs      []Input `json:"inputs"`
}

// Input type that represents input config
type Input struct {
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	Default string   `json:"default"`
	Label   string   `json:"label"`
	Items   []string `json:"items"`
	Cache   Cache    `json:"cache"`
}

type Cache struct {
	Active   bool   `json:"active"`
	Qtd      int    `json:"qtd"`
	NewLabel string `json:"newLabel"`
}

// Definition type that represents a Formula
type Definition struct {
	Path    string
	Bin     string
	Bundle	string
	Config  string
	RepoUrl string
}

// FormulaPath builds the formula path from ritchie home
func (d *Definition) FormulaPath(home string) string {
	return fmt.Sprintf(PathPattern, home, d.Path)
}

// TmpWorkDirPath builds the tmp paths to run formula, first parameter is tmpDir created
// second parameter is tmpBinDir
func (d *Definition) TmpWorkDirPath(home, uuidHash string) (string, string) {
	tmpDir := fmt.Sprintf(TmpDirPattern, home, uuidHash)
	tmpBinDir := fmt.Sprintf(TmpBinDirPattern, home, uuidHash, d.Path)
	return tmpDir, tmpBinDir
}

// BinName builds the bin name from definition params
func (d *Definition) BinName() string {
	if strings.Contains(d.Bin, "${so}") {
		so := runtime.GOOS
		suffix := ""
		if so == windows {
			suffix = ".exe"
		}
		binSO := strings.ReplaceAll(d.Bin, "${so}", so)

		return fmt.Sprintf(BinPattern, binSO, suffix)
	}
	return d.Bin
}

// BinName builds the bin name from definition params
func (d *Definition) 	BundleName() string {
	if strings.Contains(d.Bundle, "${so}") {
		so := runtime.GOOS
		bundleSO := strings.ReplaceAll(d.Bundle, "${so}", so)

		return bundleSO
	}
	return d.Bundle
}

// BinPath builds the bin path from formula path
func (d *Definition) BinPath(formula string) string {
	return fmt.Sprintf(BinPathPattern, formula)
}

// BinFilePath builds the bin file path from binPath and binName
func (d *Definition) BinFilePath(binPath, binName string) string {
	return fmt.Sprintf("%s/%s", binPath, binName)
}

// BinUrl builds the bin url
func (d *Definition) BundleUrl() string {
	return fmt.Sprintf("%s/%s/%s", d.RepoUrl, d.Path, d.BundleName())
}

// ConfigName resolver de config name
func (d *Definition) ConfigName() string {
	if d.Config != "" {
		return d.Config
	}
	return DefaultConfig
}

// ConfigPath builds the config path from formula path and config name
func (d *Definition) ConfigPath(formula, configName string) string {
	return fmt.Sprintf(ConfigPattern, formula, configName)
}

// ConfigUrl builds the config url
func (d *Definition) ConfigUrl(configName string) string {
	return fmt.Sprintf("%s/%s/%s", d.RepoUrl, d.Path, configName)
}

type Runner interface {
	Run(def Definition) error
}

type Creator interface {
	Create(formulaCmd string) error
}
