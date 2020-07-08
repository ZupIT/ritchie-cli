package creator

import (
	"errors"
	"os"
	"path"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const (
	fCmdExists        = "rit add repo"
	fCmdCorrectGo     = "rit scaffold generate test_go"
	fCmdCorrectJava   = "rit scaffold generate test_java"
	fCmdCorrectNode   = "rit scaffold generate test_node"
	fCmdCorrectPython = "rit scaffold generate test_python"
	fCmdCorrectRuby   = "rit scaffold generate test_ruby"
	fCmdCorrectShell  = "rit scaffold generate test_shell"
	fCmdCorrectPhp    = "rit scaffold generate test_php"
	langGo            = "Go"
	langJava          = "Java"
	langNode          = "Node"
	langPython        = "Python"
	langRuby          = "Ruby"
	langShell         = "Shell"
	langPhp           = "Php"
)

func TestCreator(t *testing.T) {
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	cleanForm(dirManager)

	makefileDir := createDirWithMakefile(dirManager, fileManager)
	jsonDir := createDirWithTree(dirManager, fileManager)
	fullDir := createFullDir(dirManager, fileManager)

	treeMan := tree.NewTreeManager("../../testdata", repoListerMock{}, api.SingleCoreCmds)

	type in struct {
		formCreate formula.Create
		dir        stream.DirCreater
		file       stream.FileWriteReadExister
	}

	type out struct {
		err error
	}

	tests := []struct {
		name string
		in   in
		out  out
	}{
		{
			name: "command exists",
			in: in{
				formCreate: formula.Create{
					FormulaCmd:    fCmdExists,
					Lang:          langGo,
					WorkspacePath: fullDir,
					FormulaPath:   path.Join(fullDir, "/add/repo"),
				},
				dir:  dirManager,
				file: fileManager,
			},
			out: out{
				err: ErrRepeatedCommand,
			},
		},
		{
			name: "command correct-go",
			in: in{
				formCreate: formula.Create{
					FormulaCmd:    fCmdCorrectGo,
					Lang:          langGo,
					WorkspacePath: fullDir,
					FormulaPath:   path.Join(fullDir, "/scaffold/generate/test_go"),
				},
				dir:  dirManager,
				file: fileManager,
			},
			out: out{
				err: nil,
			},
		},
		{
			name: "command correct-java",
			in: in{
				formCreate: formula.Create{
					FormulaCmd:    fCmdCorrectJava,
					Lang:          langJava,
					WorkspacePath: fullDir,
					FormulaPath:   path.Join(fullDir, "/scaffold/generate/test_java"),
				},
				dir:  dirManager,
				file: fileManager,
			},
			out: out{
				err: nil,
			},
		},
		{
			name: "command correct-node",
			in: in{
				formCreate: formula.Create{
					FormulaCmd:    fCmdCorrectNode,
					Lang:          langNode,
					WorkspacePath: fullDir,
					FormulaPath:   path.Join(fullDir, "/scaffold/generate/test_node"),
				},
				dir:  dirManager,
				file: fileManager,
			},
			out: out{
				err: nil,
			},
		},
		{
			name: "command correct-python",
			in: in{
				formCreate: formula.Create{
					FormulaCmd:    fCmdCorrectPython,
					Lang:          langPython,
					WorkspacePath: fullDir,
					FormulaPath:   path.Join(fullDir, "/scaffold/generate/test_python"),
				},
				dir:  dirManager,
				file: fileManager,
			},
			out: out{
				err: nil,
			},
		},
		{
			name: "command correct-ruby",
			in: in{
				formCreate: formula.Create{
					FormulaCmd:    fCmdCorrectRuby,
					Lang:          langRuby,
					WorkspacePath: fullDir,
					FormulaPath:   path.Join(fullDir, "/scaffold/generate/test_ruby"),
				},
				dir:  dirManager,
				file: fileManager,
			},
			out: out{
				err: nil,
			},
		},
		{
			name: "command correct-shell",
			in: in{
				formCreate: formula.Create{
					FormulaCmd:    fCmdCorrectShell,
					Lang:          langShell,
					WorkspacePath: fullDir,
					FormulaPath:   path.Join(fullDir, "/scaffold/generate/test_shell"),
				},
				dir:  dirManager,
				file: fileManager,
			},
			out: out{
				err: nil,
			},
		},
		{
			name: "command correct-php",
			in: in{
				formCreate: formula.Create{
					FormulaCmd:    fCmdCorrectPhp,
					Lang:          langPhp,
					WorkspacePath: fullDir,
					FormulaPath:   path.Join(fullDir, "/scaffold/generate/test_php"),
				},
				dir:  dirManager,
				file: fileManager,
			},
			out: out{
				err: nil,
			},
		},
		{
			name: "command to custom repo with missing tree.json",
			in: in{
				formCreate: formula.Create{
					FormulaCmd:    fCmdCorrectGo,
					Lang:          langGo,
					WorkspacePath: makefileDir,
					FormulaPath:   path.Join(makefileDir, "/scaffold/generate/test_go"),
				},
				dir:  dirManager,
				file: fileManager,
			},
			out: out{
				err: nil,
			},
		},
		{
			name: "command to custom repo with missing MakefilePath",
			in: in{
				formCreate: formula.Create{
					FormulaCmd:    fCmdCorrectGo,
					Lang:          langGo,
					WorkspacePath: jsonDir,
					FormulaPath:   path.Join(jsonDir, "/scaffold/generate/test_go"),
				},
				dir:  dirManager,
				file: fileManager,
			},
			out: out{
				err: nil,
			},
		},
		{
			name: "error create dir",
			in: in{
				formCreate: formula.Create{
					FormulaCmd:    fCmdCorrectGo,
					Lang:          langGo,
					WorkspacePath: fullDir,
					FormulaPath:   path.Join(jsonDir, "/scaffold/generate/test_go"),
				},
				dir:  dirManagerMock{createErr: errors.New("error to create dir")},
				file: fileManager,
			},
			out: out{
				err: errors.New("error to create dir"),
			},
		},
		{
			name: "error write file",
			in: in{
				formCreate: formula.Create{
					FormulaCmd:    fCmdCorrectGo,
					Lang:          langGo,
					WorkspacePath: fullDir,
					FormulaPath:   path.Join(jsonDir, "/scaffold/generate/test_go"),
				},
				dir:  dirManager,
				file: fileManagerMock{writeErr: errors.New("error to write file")},
			},
			out: out{
				err: errors.New("error to write file"),
			},
		},
		{
			name: "error read file",
			in: in{
				formCreate: formula.Create{
					FormulaCmd:    "rit test error",
					Lang:          langGo,
					WorkspacePath: fullDir,
					FormulaPath:   path.Join(fullDir, "/test/error"),
				},
				dir:  dirManager,
				file: fileManagerMock{readErr: errors.New("error to read file"), exist: true},
			},
			out: out{
				err: errors.New("error to read file"),
			},
		},
		{
			name: "error read json",
			in: in{
				formCreate: formula.Create{
					FormulaCmd:    "rit test error",
					Lang:          langGo,
					WorkspacePath: fullDir,
					FormulaPath:   path.Join(fullDir, "/test/error"),
				},
				dir:  dirManager,
				file: fileManagerMock{data: []byte(""), exist: true},
			},
			out: out{
				err: errors.New("unexpected end of JSON input"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in
			creator := NewCreator(treeMan, tt.in.dir, tt.in.file)
			out := tt.out
			got := creator.Create(in.formCreate)
			if got != nil && got.Error() != out.err.Error() || out.err != nil && got == nil {
				t.Errorf("Create(%s) got %v, want %v", tt.name, got, out.err)
			}
		})
	}
}

func TestCreatorFail(t *testing.T) {
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	cleanForm(dirManager)

	fullDir := createFullDir(dirManager, fileManager)

	treeMan := tree.NewTreeManager("../../testdata", repoListerMock{}, api.SingleCoreCmds)

	tests := []string{langGo, langJava, langNode, langPhp, langPython, langShell}

	creatorMock := genericFileCreatorMock{ createErr: errors.New("error while creating language") }
	creator := NewCreator(treeMan, dirManager, fileManager)
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			expected := errors.New("error while creating language")
			formulaPath := path.Join(fullDir, "/scaffold/generate/test_fail")
			got := creator.createSrcFiles(formulaPath, "test_fail", tt, creatorMock)
			if got == nil || got.Error() != expected.Error() {
				t.Errorf("Create Formula Fail(%s) got %v, want %v", tt, got, expected)
			}
		})
	}
}

type repoListerMock struct{}

func (repoListerMock) List() ([]formula.Repository, error) {
	return []formula.Repository{}, nil
}

type dirManagerMock struct {
	createErr error
}

func (d dirManagerMock) Create(string) error {
	return d.createErr
}

type fileManagerMock struct {
	data     []byte
	writeErr error
	readErr  error
	exist    bool
}

func (f fileManagerMock) Write(string, []byte) error {
	return f.writeErr
}
func (f fileManagerMock) Read(string) ([]byte, error) {
	return f.data, f.readErr
}
func (f fileManagerMock) Exists(string) bool {
	return f.exist
}

type genericFileCreatorMock struct {
	createErr   error
	GenericFileCreatorI
}

func (c genericFileCreatorMock) createGenericFiles(_, _, _ string, _ formula.Lang) error {
	return c.createErr
}

func cleanForm(dir stream.DirManager) {
	tempDir := os.TempDir()
	_ = dir.Remove(path.Join(tempDir, formula.DefaultWorkspaceDir))
	_ = dir.Remove(path.Join(tempDir, "/customRepo"))
	_ = dir.Remove(path.Join(tempDir, "/customRepoMakefile"))
	_ = dir.Remove(path.Join(tempDir, "/customRepoTreejson"))
}

func createDirWithMakefile(dir stream.DirCreater, file stream.FileWriter) string {
	workspacePath := path.Join(os.TempDir(), "/customRepoMakefile")
	makefilePath := path.Join(workspacePath, formula.MakefilePath)
	_ = dir.Create(workspacePath)
	_ = file.Write(makefilePath, []byte(""))
	return workspacePath
}

func createDirWithTree(dir stream.DirCreater, file stream.FileWriter) string {
	workspacePath := path.Join(os.TempDir(), "/customRepoTreejson")
	treeJsonDir := path.Join(workspacePath, "/tree")
	treeJsonFile := path.Join(workspacePath, formula.TreePath)
	_ = dir.Create(workspacePath)
	_ = dir.Create(treeJsonDir)
	_ = file.Write(treeJsonFile, []byte("{}"))
	return workspacePath
}

func createFullDir(dir stream.DirCreater, file stream.FileWriteReadExister) string {
	workspacePath := path.Join(os.TempDir(), "/customRepo")
	treeJsonDir := path.Join(workspacePath, "/tree")
	treeJsonFile := path.Join(workspacePath, formula.TreePath)
	makefilePath := path.Join(workspacePath, formula.MakefilePath)
	_ = dir.Create(workspacePath)
	_ = dir.Create(treeJsonDir)
	makefile, _ := file.Read("../../testdata/MakefilePath")
	_ = file.Write(makefilePath, makefile)
	_ = file.Write(treeJsonFile, []byte("{}"))

	return workspacePath
}
