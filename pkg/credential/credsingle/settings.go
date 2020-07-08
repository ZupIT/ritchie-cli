package credsingle

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type SingleSettings struct {
	file stream.FileWriteReadExister
}

func NewSingleSettings(file stream.FileWriteReadExister) SingleSettings {
	return SingleSettings{file: file}
}

func (s SingleSettings) ReadCredentials(path string) (credential.Fields, error) {

	var fields credential.Fields

	if s.file.Exists(path) {
		cBytes, _ := s.file.Read(path)
		err := json.Unmarshal(cBytes, &fields)
		if err != nil {
			return fields, err
		}
	}

	return fields, nil
}

func (s SingleSettings) WriteCredentials(fields credential.Fields, path string) error {
	fieldsData, err := json.Marshal(fields)
	if err != nil {
		return err
	}
	err = s.file.Write(path, fieldsData)
	if err != nil {
		return err
	}

	return nil
}

func (s SingleSettings) DefaultCredentials() error {
	if !s.file.Exists(ProviderPath()) {
		err := s.WriteCredentials(NewDefaultCredentials(), ProviderPath())
		return err
	}
	return nil
}

func ProviderPath() string {
	homeDir, _ := os.UserHomeDir()
	providerDir := fmt.Sprintf("%s/.rit/repo/providers.json", homeDir)
	return providerDir
}

func NewDefaultCredentials() credential.Fields {

	var username = credential.Field{
		Name: "username",
		Type: "text",
	}

	var token = credential.Field{
		Name: "token",
		Type: "password",
	}

	var accessKeyId = credential.Field{
		Name: "accessKeyId",
		Type: "text",
	}

	var secretAccessKey = credential.Field{
		Name: "secretAccessKey",
		Type: "password",
	}

	var base64config = credential.Field{
		Name: "base64config",
		Type: "text",
	}

	var dc = credential.Fields{
		"Add a new":  []credential.Field{},
		"github":     []credential.Field{username, token},
		"gitlab":     []credential.Field{username, token},
		"aws":        []credential.Field{accessKeyId, secretAccessKey},
		"jenkins":    []credential.Field{username, token},
		"kubeconfig": []credential.Field{base64config},
	}

	return dc
}
