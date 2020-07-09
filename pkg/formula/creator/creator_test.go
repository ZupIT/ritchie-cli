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
	fCmdCorrectShell  = "rit scaffold generate test_shell"
	fCmdCorrectPhp    = "rit scaffold generate test_php"
	langGo            = "Go"
	langJava          = "Java"
	langNode          = "Node"
	langPython        = "Python"
	langShell         = "Shell"
	langPhp           = "Php"
)

func TestCreator(t *testing.T) {
	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	cleanForm(dirManager)

	gitIgnoreDir := createDirWithGitIgnore(dirManager, fileManager)
	mainReadMeDir := createDirMainReadMe(dirManager, fileManager)
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
			name: "command to custom repo with missing ReadMe",
			in: in{
				formCreate: formula.Create{
					FormulaCmd:    fCmdCorrectGo,
					Lang:          langGo,
					WorkspacePath: gitIgnoreDir,
					FormulaPath:   path.Join(gitIgnoreDir, "/scaffold/generate/test_go"),
				},
				dir:  dirManager,
				file: fileManager,
			},
			out: out{
				err: nil,
			},
		},
		{
			name: "command to custom repo with missing GitIgnore",
			in: in{
				formCreate: formula.Create{
					FormulaCmd:    fCmdCorrectGo,
					Lang:          langGo,
					WorkspacePath: mainReadMeDir,
					FormulaPath:   path.Join(mainReadMeDir, "/scaffold/generate/test_go"),
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
					FormulaPath:   path.Join(mainReadMeDir, "/scaffold/generate/test_go"),
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
					FormulaPath:   path.Join(mainReadMeDir, "/scaffold/generate/test_go"),
				},
				dir:  dirManager,
				file: fileManagerMock{writeErr: errors.New("error to write file")},
			},
			out: out{
				err: errors.New("error to write file"),
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

func TestCreateManager_existsGitIgnore(t *testing.T) {

	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	cleanForm(dirManager)

	treeMan := tree.NewTreeManager("../../testdata", repoListerMock{}, api.SingleCoreCmds)

	type fields struct {
		treeManager tree.Manager
		dir         stream.DirCreater
		file        stream.FileWriteReadExister
	}
	type args struct {
		workspacePath string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: ".gitignore exist",
			fields: fields{
				treeManager: treeMan,
				dir:         dirManager,
				file: fileManagerMock{
					data:    []byte("some text"),
					readErr: nil,
					exist:   true,
				},
			},
			args: args{workspacePath: ""},
			want: true,
		},
		{
			name: ".gitignore not exist",
			fields: fields{
				treeManager: treeMan,
				dir:         dirManager,
				file: fileManagerMock{
					data:    []byte("some text"),
					readErr: nil,
					exist:   false,
				},
			},
			args: args{workspacePath: ""},
			want: false,
		},
		{
			name: ".gitignore err to read",
			fields: fields{
				treeManager: treeMan,
				dir:         dirManager,
				file: fileManagerMock{
					data:    []byte(""),
					readErr: errors.New("some errors"),
					exist:   true,
				},
			},
			args: args{workspacePath: ""},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CreateManager{
				treeManager: tt.fields.treeManager,
				dir:         tt.fields.dir,
				file:        tt.fields.file,
			}
			if got := c.existsGitIgnore(tt.args.workspacePath); got != tt.want {
				t.Errorf("existsGitIgnore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateManager_existsMainReadMe(t *testing.T) {

	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)
	cleanForm(dirManager)

	treeMan := tree.NewTreeManager("../../testdata", repoListerMock{}, api.SingleCoreCmds)

	type fields struct {
		treeManager tree.Manager
		dir         stream.DirCreater
		file        stream.FileWriteReadExister
	}
	type args struct {
		workspacePath string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "readMe exist",
			fields: fields{
				treeManager: treeMan,
				dir:         dirManager,
				file: fileManagerMock{
					data:    []byte("some text"),
					readErr: nil,
					exist:   true,
				},
			},
			args: args{workspacePath: ""},
			want: true,
		},
		{
			name: "readMe not exist",
			fields: fields{
				treeManager: treeMan,
				dir:         dirManager,
				file: fileManagerMock{
					data:    []byte("some text"),
					readErr: nil,
					exist:   false,
				},
			},
			args: args{workspacePath: ""},
			want: false,
		},
		{
			name: "readMe fail to read",
			fields: fields{
				treeManager: treeMan,
				dir:         dirManager,
				file: fileManagerMock{
					data:    []byte("some text"),
					readErr: errors.New("some error"),
					exist:   true,
				},
			},
			args: args{workspacePath: ""},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CreateManager{
				treeManager: tt.fields.treeManager,
				dir:         tt.fields.dir,
				file:        tt.fields.file,
			}
			if got := c.existsMainReadMe(tt.args.workspacePath); got != tt.want {
				t.Errorf("existsMainReadMe() = %v, want %v", got, tt.want)
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

func cleanForm(dir stream.DirManager) {
	tempDir := os.TempDir()
	_ = dir.Remove(path.Join(tempDir, formula.DefaultWorkspaceDir))
	_ = dir.Remove(path.Join(tempDir, "/customRepo"))
	_ = dir.Remove(path.Join(tempDir, "/customRepoMakefile"))
	_ = dir.Remove(path.Join(tempDir, "/customRepoTreejson"))
}

func createDirWithGitIgnore(dir stream.DirCreater, file stream.FileWriter) string {
	workspacePath := path.Join(os.TempDir(), "/customRepoGitIgnore")
	gitIgnorePath := path.Join(workspacePath, ".gitignore")
	_ = dir.Create(workspacePath)
	_ = file.Write(gitIgnorePath, []byte(""))
	return workspacePath
}

func createDirMainReadMe(dir stream.DirCreater, file stream.FileWriter) string {
	workspacePath := path.Join(os.TempDir(), "/customRepoReadMe")
	readMePath := path.Join(workspacePath, "README.md")
	_ = dir.Create(workspacePath)
	_ = file.Write(readMePath, []byte("{}"))
	return workspacePath
}

func createFullDir(dir stream.DirCreater, file stream.FileWriteReadExister) string {
	workspacePath := path.Join(os.TempDir(), "/customRepo")
	gitIgnorePath := path.Join(workspacePath, ".gitignore")
	readMePath := path.Join(workspacePath, "README.md")
	_ = dir.Create(workspacePath)
	_ = file.Write(gitIgnorePath, []byte("{}"))
	_ = file.Write(readMePath, []byte("{}"))

	return workspacePath
}
