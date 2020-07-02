package runner

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

func TestOutputManager_ValidAndPrint(t *testing.T) {

	tmpDir := os.TempDir() + "/Test_printAndValidOutputDir"
	_ = fileutil.CreateDirIfNotExists(tmpDir, 0755)
	defer func() { _ = fileutil.RemoveDir(tmpDir) }()

	type args struct {
		setup formula.Setup
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		err     error
	}{
		{
			name: "Return empty string when dir is empty",
			args: args{
				setup: formula.Setup{
					Config: formula.Config{Outputs: []formula.Output{}},
					TmpOutputDir: func() string {
						basePath := "/t-rit-return-empty"
						path := tmpDir + basePath
						_ = fileutil.CreateDirIfNotExists(path, 0755)
						return path
					}(),
				},
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "Return only the outputs with printValue",
			args: args{
				setup: formula.Setup{
					Config: formula.Config{Outputs: []formula.Output{
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
					TmpOutputDir: func() string {
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
			want:    "X=1\nZ=3\n",
			wantErr: false,
		},
		{
			name: "Return Err when output dir not have all files",
			args: args{
				setup: formula.Setup{
					Config: formula.Config{Outputs: []formula.Output{
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
					TmpOutputDir: func() string {
						basePath := "/t-rit-err-all-files"
						path := tmpDir + basePath
						_ = fileutil.CreateDirIfNotExists(path, 0755)
						_ = ioutil.WriteFile(path+"/x", []byte("1"), 0755)
						_ = ioutil.WriteFile(path+"/z", []byte("3"), 0755)
						return path
					}(),
				},
			},
			want:    "",
			wantErr: true,
			err:     ErrValidOutputDir,
		},
		{
			name: "Return when some output file is missing",
			args: args{
				setup: formula.Setup{
					Config: formula.Config{Outputs: []formula.Output{
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
					TmpOutputDir: func() string {
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
			want:    "",
			wantErr: true,
			err:     errors.New(prompt.Red("file:Y not found in output dir")),
		},
		{
			name: "Return Err when fail to read dir",
			args: args{
				setup: formula.Setup{
					Config: formula.Config{Outputs: []formula.Output{}},
					TmpOutputDir: func() string {
						basePath := "/not-created-dir"
						return basePath
					}(),
				},
			},
			want:    "",
			wantErr: true,
			err:     ErrReadOutputDir,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := bytes.Buffer{}
			o := NewOutputManager(&buffer)

			err := o.Outputs(tt.args.setup)
			if (err != nil) != tt.wantErr {
				t.Errorf("Outputs(%s) error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if got := buffer.String(); got != tt.want {
				t.Errorf("Outputs(%s) = %v, want %v", tt.name, got, tt.want)
			}
			if err != nil && tt.wantErr && err.Error() != tt.err.Error() {
				t.Errorf("Outputs(%s) = err:%v, wantErr %v", tt.name, err, tt.err)
			}
		})
	}
}
