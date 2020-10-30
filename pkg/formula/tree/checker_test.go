package tree

import (
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
	sMocks "github.com/ZupIT/ritchie-cli/pkg/stream/mocks"
)

func TestChecker(t *testing.T) {
	treeJson := `{
					"commands": [
						{
							"id": "root_aws_create_bucket",
							"parent": "root_aws_create",
							"usage": "bucket",
							"help": "short help placeholder for bucket",
							"longHelp": "long help placeholder for bucket used by index page and -h",
							"formula": true
						},
						{
							"id": "root_aws_create_bucket",
							"parent": "root_aws_create",
							"usage": "bucket",
							"help": "short help placeholder for bucket",
							"longHelp": "long help placeholder for bucket used by index page and -h",
							"formula": true
						}
					]
				}
`

	tests := []struct{
		name string
		file stream.FileReader
		dir DirManagerCustomMock
	}{
		{
			name: "Should success run",
			file : sMocks.FileReaderCustomMock{
				ReadMock: func(path string) ([]byte, error) {
					return []byte(treeJson), nil
				},
			},
			dir: DirManagerCustomMock{
				list: func(dir string, hiddenDir bool) ([]string, error) {
					return []string{"commons"}, nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			treeChecker := NewChecker(tt.dir, tt.file)
			treeChecker.CheckCommands()
		})
	}

}
type DirManagerCustomMock struct {
	exists func(dir string) bool
	list   func(dir string, hiddenDir bool) ([]string, error)
	isDir  func(dir string) bool
	create func(dir string) error
}

func (d DirManagerCustomMock) Exists(dir string) bool {
	return d.exists(dir)
}

func (d DirManagerCustomMock) List(dir string, hiddenDir bool) ([]string, error) {
	return d.list(dir, hiddenDir)
}

func (d DirManagerCustomMock) IsDir(dir string) bool {
	return d.isDir(dir)
}

func (d DirManagerCustomMock) Create(dir string) error {
	return d.create(dir)
}