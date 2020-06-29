package formula

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

func TestDefaultRunner_Run(t *testing.T) {
	def := Definition{
		Path:    "mock/test",
		Bin:     "test-linux",
		Bundle:  "linux.zip",
		Config:  "config.json",
		RepoURL: RepoUrl,
	}

	home := os.TempDir()
	_ = fileutil.RemoveDir(home + "/formulas")
	setup := NewDefaultSingleSetup(home, http.DefaultClient)

	type in struct {
		envMock  envResolverMock
		inText   inputMock
		inBool   inputMock
		inPass   inputMock
		preMock  *preRunnerMock
		postMock *postRunnerMock
	}

	tests := []struct {
		name string
		in   in
		want error
	}{
		{
			name: "success",
			in: in{
				envMock: envResolverMock{in: "ok"},
				inText:  inputMock{text: ""},
				inBool:  inputMock{boolean: true},
				inPass:  inputMock{text: "******"},
			},
			want: nil,
		},
		{
			name: "pre run error",
			in: in{
				envMock: envResolverMock{in: "ok"},
				inText:  inputMock{text: "ok"},
				inBool:  inputMock{boolean: true},
				inPass:  inputMock{text: "******"},
				preMock: &preRunnerMock{
					setup: Setup{},
					error: ErrFormulaBinNotFound,
				},
			},
			want: ErrFormulaBinNotFound,
		},
		{
			name: "inputs error",
			in: in{
				envMock: envResolverMock{in: "ok"},
				inText:  inputMock{err: errors.New("fail to resolve input")},
				inBool:  inputMock{boolean: true},
				inPass:  inputMock{text: "******"},
			},
			want: errors.New("fail to resolve input"),
		},
		{
			name: "post run error",
			in: in{
				envMock:  envResolverMock{in: "ok"},
				inText:   inputMock{text: "ok"},
				inBool:   inputMock{boolean: true},
				inPass:   inputMock{text: "******"},
				postMock: &postRunnerMock{error: errors.New("error in remove dir")},
			},
			want: errors.New("error in remove dir"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in

			var preRunner PreRunner
			if in.preMock != nil {
				preRunner = in.preMock
			} else {
				preRunner = NewDefaultPreRunner(setup)
			}

			var postRunner PostRunner
			if in.postMock != nil {
				postRunner = in.postMock
			} else {
				postRunner = NewPostRunner()
			}

			resolvers := env.Resolvers{"test": in.envMock}
			inputManager := NewInputManager(resolvers, in.inText, in.inText, in.inBool, in.inPass)
			defaultRunner := NewDefaultRunner(preRunner, postRunner, inputManager)

			got := defaultRunner.Run(def, api.Prompt)

			if got != nil && got.Error() != tt.want.Error() {
				t.Errorf("Run(%s) got %v, want %v", tt.name, got, tt.want)
			}

		})
	}
}

func Test_printAndValidOutputDir(t *testing.T) {

	tmpDir := os.TempDir() + "/Test_printAndValidOutputDir"
	_ = fileutil.CreateDirIfNotExists(tmpDir, 0755)
	defer func() { _ = fileutil.RemoveDir(tmpDir) }()

	type args struct {
		setup Setup
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Return empty string when dir is empty",
			args: args{
				setup: Setup{
					config: Config{Outputs: []Output{}},
					tmpOutputDir: func() string {
						basePath := "/t-rit-return-empty"
						path := tmpDir + basePath
						_ = fileutil.CreateDirIfNotExists(path, 0755)
						return path
					}(),
				},
			},
			want: "",
		},
		{
			name: "Return only the outputs with printValue",
			args: args{
				setup: Setup{
					config: Config{Outputs: []Output{
						{
							Name:  "X",
							Print: true,
						},
						{
							Name:  "Y",
							Print: false,
						},
						{
							Name:  "Z",
							Print: true,
						},
					}},
					tmpOutputDir: func() string {
						basePath := "/t-rit-printed"
						path := tmpDir + basePath
						_ = fileutil.CreateDirIfNotExists(path, 0755)
						_ = ioutil.WriteFile(path+"/x", []byte("1"), 0755)
						_ = ioutil.WriteFile(path+"/y", []byte("2"), 0755)
						_ = ioutil.WriteFile(path+"/z", []byte("3"), 0755)
						return path
					}(),
				},
			},
			want: "X=1\nZ=3\n",
		},
		{
			name: "Return Red when output dir not have all files",
			args: args{
				setup: Setup{
					config: Config{Outputs: []Output{
						{
							Name:  "X",
							Print: true,
						},
						{
							Name:  "Y",
							Print: false,
						},
						{
							Name:  "Z",
							Print: true,
						},
					}},
					tmpOutputDir: func() string {
						basePath := "/t-rit-err-all-files"
						path := tmpDir + basePath
						_ = fileutil.CreateDirIfNotExists(path, 0755)
						_ = ioutil.WriteFile(path+"/x", []byte("1"), 0755)
						_ = ioutil.WriteFile(path+"/z", []byte("3"), 0755)
						return path
					}(),
				},
			},
			want: prompt.Red("Output dir not have all the outputs files"),
		},
		{
			name: "Return Red when some output file is missing",
			args: args{
				setup: Setup{
					config: Config{Outputs: []Output{
						{
							Name:  "X",
							Print: true,
						},
						{
							Name:  "Y",
							Print: false,
						},
						{
							Name:  "Z",
							Print: true,
						},
					}},
					tmpOutputDir: func() string {
						basePath := "/t-rit-err-missing-files"
						path := tmpDir + basePath
						_ = fileutil.CreateDirIfNotExists(path, 0755)
						_ = ioutil.WriteFile(path+"/x", []byte("1"), 0755)
						_ = ioutil.WriteFile(path+"/z", []byte("3"), 0755)
						_ = ioutil.WriteFile(path+"/w", []byte("3"), 0755)
						return path
					}(),
				},
			},
			want: prompt.Red("file:Y not found in output dir"),
		},
		{
			name: "Return Err when fail to read dir",
			args: args{
				setup: Setup{
					config: Config{Outputs: []Output{}},
					tmpOutputDir: func() string {
						basePath := "/not-created-dir"
						return basePath
					}(),
				},
			},
			want: prompt.Red("Fail to read output dir"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := printAndValidOutputDir(tt.args.setup); got != tt.want {
				t.Errorf("printAndValidOutputDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

type inputMock struct {
	text    string
	boolean bool
	err     error
}

func (i inputMock) List(string, []string) (string, error) {
	return i.text, i.err
}

func (i inputMock) Text(string, bool) (string, error) {
	return i.text, i.err
}

func (i inputMock) Bool(string, []string) (bool, error) {
	return i.boolean, i.err
}

func (i inputMock) Password(string) (string, error) {
	return i.text, i.err
}

type envResolverMock struct {
	in  string
	err error
}

func (e envResolverMock) Resolve(string) (string, error) {
	return e.in, e.err
}

type preRunnerMock struct {
	setup Setup
	error error
}

func (p preRunnerMock) PreRun(Definition) (Setup, error) {
	return p.setup, p.error
}

type postRunnerMock struct {
	error error
}

func (p postRunnerMock) PostRun(Setup, bool) error {
	return p.error
}
