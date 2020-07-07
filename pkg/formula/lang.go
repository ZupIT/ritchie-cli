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
