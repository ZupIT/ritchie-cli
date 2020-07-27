package repo

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func TestSetPriorityManager_SetPriority(t *testing.T) {

	fileManager := stream.NewFileManager()
	dirManager := stream.NewDirManager(fileManager)

	type fields struct {
		ritHome string
		file    stream.FileWriteReadExister
	}
	type args struct {
		repoName formula.RepoName
		priority int
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		Err     error
	}{
		{
			name: "Setting priority test success",
			fields: fields{
				ritHome: func() string {
					ritHomePath := filepath.Join(os.TempDir(), "test-priority-setter-repo-sucess")
					_ = dirManager.Remove(ritHomePath)
					_ = dirManager.Create(ritHomePath)
					_ = dirManager.Create(filepath.Join(ritHomePath, "repos"))

					repositoryFile := filepath.Join(ritHomePath, "repos", "repositories.json")

					data := `
						[
							{
								"name": "commons",
								"version": "v2.0.0",
								"url": "https://github.com/kaduartur/ritchie-formulas",
								"priority": 0
							}
						]`

					_ = fileManager.Write(repositoryFile, []byte(data))
					return ritHomePath
				}(),
				file: fileManager,
			},
			args: args{
				repoName: "commons",
				priority: 1,
			},
			wantErr: false,
		},
		{
			name: "Return error when try to unmarshal the file to json",
			fields: fields{
				ritHome: func() string {
					ritHomePath := filepath.Join(os.TempDir(), "test-priority-setter-repo-fail")
					_ = dirManager.Remove(ritHomePath)
					_ = dirManager.Create(ritHomePath)
					_ = dirManager.Create(filepath.Join(ritHomePath, "repos"))

					repositoryFile := filepath.Join(ritHomePath, "repos", "repositories.json")

					data := `
						[
							{
								"errorHere: "commons",
								"version": "v2.0.0",
								"url": "https://github.com/kaduartur/ritchie-formulas",
								"priority": 0
							}
						]`

					_ = fileManager.Write(repositoryFile, []byte(data))
					return ritHomePath
				}(),
				file: fileManager,
			},
			args: args{
				repoName: "commons",
				priority: 1,
			},
			wantErr: true,
		},
		{
			name: "Return error when file not exist",
			fields: fields{
				ritHome: os.TempDir(),
				file:    fileManager,
			},
			args: args{
				repoName: "commons",
				priority: 1,
			},
			wantErr: true,
			Err:     errors.New(repositoryDoNotExistError),
		},
		{
			name: "Return error when try to read file",
			fields: fields{
				ritHome: os.TempDir(),
				file:    fileWriteReadExisterMockErrorOnReadAndWrite{},
			},
			args: args{
				repoName: "commons",
				priority: 1,
			},
			wantErr: true,
		},
		{
			name: "Return error when try to write the changes on file",
			fields: fields{
				ritHome: os.TempDir(),
				file:    fileWriteReadExisterMockOnSucessReadData{},
			},
			args: args{
				repoName: "commons",
				priority: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			sm := SetPriorityManager{
				ritHome: tt.fields.ritHome,
				file:    tt.fields.file,
			}

			err := sm.SetPriority(tt.args.repoName, tt.args.priority)

			if (tt.Err != nil) && err.Error() != tt.Err.Error() {
				t.Errorf("This error didnt expect this menssage")
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("SetPriorityManager.SetPriority() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type fileWriteReadExisterMockOnSucessReadData struct{}

func (m fileWriteReadExisterMockOnSucessReadData) Read(path string) ([]byte, error) {
	dataWithoutErrors := `
	[
		{
			"name": "commons",
			"version": "v2.0.0",
			"url": "https://github.com/kaduartur/ritchie-formulas",
			"priority": 0
		}
	]`
	return []byte(dataWithoutErrors), nil
}

func (m fileWriteReadExisterMockOnSucessReadData) Write(path string, content []byte) error {
	return errors.New("Error on write the data on file")
}

func (m fileWriteReadExisterMockOnSucessReadData) Exists(path string) bool {
	return true
}

type fileWriteReadExisterMockErrorOnReadAndWrite struct{}

func (m fileWriteReadExisterMockErrorOnReadAndWrite) Read(path string) ([]byte, error) {
	return nil, errors.New("Error on read the file")
}

func (m fileWriteReadExisterMockErrorOnReadAndWrite) Exists(path string) bool {
	return true
}

func (m fileWriteReadExisterMockErrorOnReadAndWrite) Write(path string, content []byte) error {
	return errors.New("Error on write the data on file")
}
