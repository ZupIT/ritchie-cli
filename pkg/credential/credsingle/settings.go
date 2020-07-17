package credsingle

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const AddNew = "Add a new"

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

// WriteDefault is a non override version of WriteCredentials
// used to create providers.json if user dont have it
func (s SingleSettings) WriteDefaultCredentials(path string) error {
	if !s.file.Exists(path) {
		err := s.WriteCredentials(NewDefaultCredentials(), path)
		return err
	}
	return nil
}

func NewDefaultCredentials() credential.Fields {
	var username = credential.Field{
		Name: "username",
		Type: "plain text",
	}

	var token = credential.Field{
		Name: "token",
		Type: "secret",
	}

	var accessKeyId = credential.Field{
		Name: "accessKeyId",
		Type: "plain text",
	}

	var secretAccessKey = credential.Field{
		Name: "secretAccessKey",
		Type: "secret",
	}

	var base64config = credential.Field{
		Name: "base64config",
		Type: "plain text",
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

func ProviderPath() string {
	homeDir, _ := os.UserHomeDir()
	providerDir := fmt.Sprintf("%s/.rit/repo/providers.json", homeDir)
	return providerDir
}

func NewProviderArr(fields credential.Fields) []string {
	var providerArr []string
	for k := range fields {
		if k != AddNew {
			providerArr = append(providerArr, k)
		}
	}
	providerArr = append(providerArr, AddNew)
	return providerArr
}
