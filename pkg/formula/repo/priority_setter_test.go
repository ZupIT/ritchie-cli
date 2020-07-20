package repo

import (
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
		dir     stream.DirCreater
	}
	type args struct {
		repoName formula.RepoName
		priority int
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    formula.Repos
		wantErr bool
	}{
		{
			name: "Setting priority test success",
			fields: fields{
				ritHome: func() string {
					ritHomePath := filepath.Join(os.TempDir(), "test-priority-setter-repo")
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
				dir:  dirManager,
			},
			args: args{
				repoName: "commons",
				priority: 0,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := SetPriorityManager{
				ritHome: tt.fields.ritHome,
				file:    tt.fields.file,
				dir:     tt.fields.dir,
			}
			if err := sm.SetPriority(tt.args.repoName, tt.args.priority); (err != nil) != tt.wantErr {
				t.Errorf("SetPriorityManager.SetPriority() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
