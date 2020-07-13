package creator

import (
	"os"
	"path"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/template"
	"github.com/ZupIT/ritchie-cli/pkg/formula/tree"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const (
	fCmdExists        = "rit add repo"
	fCmdCorrectGo     = "rit scaffold generate test_go"
	fCmdCorrectJava   = "rit scaffold generate test_java"
	fCmdCorrectNode   = "rit scaffold generate test_node"
	fCmdCorrectPython = "rit scaffold generate test_python"
	fCmdCorrectShell  = "rit scaffold generate test_shell"
	langGo            = "go"
	langJava          = "java"
	langNode          = "node"
	langPython        = "python"
	langShell         = "bash shell"
)

func TestCreator(t *testing.T) {
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)

	api.RitchieHomeDir()

	resultDir := path.Join(os.TempDir(), "/customWorkSpace")
	_ = dirManager.Remove(resultDir)
	_ = dirManager.Create(resultDir)

	treeMan := tree.NewTreeManager("../../testdata", repoListerMock{}, api.CoreCmds)

	tplM := template.NewManager("../../../testdata")

	type in struct {
		formCreate formula.Create
		dir        stream.DirCreater
		file       stream.FileWriteReadExister
		tplM       template.Manager
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
					WorkspacePath: resultDir,
					FormulaPath:   path.Join(resultDir, "/add/repo"),
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
					WorkspacePath: resultDir,
					FormulaPath:   path.Join(resultDir, "/scaffold/generate/test_go"),
				},
				dir:  dirManager,
				file: fileManager,
				tplM: tplM,
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
					WorkspacePath: resultDir,
					FormulaPath:   path.Join(resultDir, "/scaffold/generate/test_java"),
				},
				dir:  dirManager,
				file: fileManager,
				tplM: tplM,
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
					WorkspacePath: resultDir,
					FormulaPath:   path.Join(resultDir, "/scaffold/generate/test_node"),
				},
				dir:  dirManager,
				file: fileManager,
				tplM: tplM,
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
					WorkspacePath: resultDir,
					FormulaPath:   path.Join(resultDir, "/scaffold/generate/test_python"),
				},
				dir:  dirManager,
				file: fileManager,
				tplM: tplM,
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
					WorkspacePath: resultDir,
					FormulaPath:   path.Join(resultDir, "/scaffold/generate/test_shell"),
				},
				dir:  dirManager,
				file: fileManager,
				tplM: tplM,
			},
			out: out{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in
			creator := NewCreator(treeMan, tt.in.dir, tt.in.file, tt.in.tplM)
			out := tt.out
			got := creator.Create(in.formCreate)
			if (got != nil && out.err == nil) || got != nil && got.Error() != out.err.Error() || out.err != nil && got == nil {
				t.Errorf("Create(%s) got %v, want %v", tt.name, got, out.err)
			}
		})
	}
}

type repoListerMock struct{}

func (repoListerMock) List() (formula.Repos, error) {
	return formula.Repos{}, nil
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
