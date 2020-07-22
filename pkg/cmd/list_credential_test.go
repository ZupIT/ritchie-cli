package cmd

import (
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

func Test_ListCredentialCmd(t *testing.T) {
	fileManager := stream.NewFileManager()
	credSettings := credential.NewSettings(fileManager)

	t.Run("Success case", func(t *testing.T) {
		o := NewListCredentialCmd(credSettings)
		if err := o.Execute();err !=nil{
			t.Errorf("Test_ListCredentialCmd error = %s", err)
		}
	})

}
