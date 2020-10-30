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
						}
					]
				}
`

	tests := []struct{
		name string
		file stream.FileReader
	}{
		{
			name: "Should success run",
			file : sMocks.FileReaderCustomMock{
				ReadMock: func(path string) ([]byte, error) {
					return []byte(treeJson), nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			treeChecker := NewChecker(stream.DirManager{}, tt.file)
			treeChecker.CheckCommands()
		})
	}

}
