package credential

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const AddNew = "Add a new"

type Settings struct {
	file stream.FileWriteReadExisterLister
	dir  stream.DirLister
}

func NewSettings(file stream.FileWriteReadExisterLister, dir stream.DirLister) Settings {
	return Settings{
		file: file,
		dir:  dir,
	}
}

func (s Settings) ReadCredentialsFields(path string) (Fields, error) {
	fields := Fields{}
	if s.file.Exists(path) {
		cBytes, _ := s.file.Read(path)
		if err := json.Unmarshal(cBytes, &fields); err != nil {
			return fields, err
		}
	}
	return fields, nil
}

func (s Settings) ReadCredentialsValue(path string) ([]ListCredData, error) {
	var creds []ListCredData
	var cred ListCredData
	var detail Detail
	ctx, _ := s.dir.List(path, true)
	for _, c := range ctx {
		providers, _ := s.file.List(filepath.Join(path, c))
		for _, p := range providers {
			cBytes, _ := s.file.Read(filepath.Join(path, c, p))
			if err := json.Unmarshal(cBytes, &detail); err != nil {
				return creds, err
			}
			for k, v := range detail.Credential {
				cred.Provider = detail.Service
				cred.Context = c
				cred.Value = v
				cred.Name = k
				creds = append(creds, cred)
			}
		}
	}
	return creds, nil
}

func (s Settings) WriteCredentialsFields(fields Fields, path string) error {
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

// WriteDefault is a non override version of WriteCredentialsFields
// used to create providers.json if user dont have it
func (s Settings) WriteDefaultCredentialsFields(path string) error {
	if !s.file.Exists(path) {
		err := s.WriteCredentialsFields(NewDefaultCredentials(), path)
		return err
	}
	return nil
}

func NewDefaultCredentials() Fields {
	username := Field{
		Name: "username",
		Type: "text",
	}

	token := Field{
		Name: "token",
		Type: "password",
	}

	accessKey := Field{
		Name: "accessKey",
		Type: "text",
	}

	secretAccessKey := Field{
		Name: "secretAccessKey",
		Type: "password",
	}

	base64config := Field{
		Name: "base64config",
		Type: "text",
	}

	dc := Fields{
		AddNew:       []Field{},
		"github":     []Field{username, token},
		"gitlab":     []Field{username, token},
		"aws":        []Field{accessKey, secretAccessKey},
		"jenkins":    []Field{username, token},
		"kubeconfig": []Field{base64config},
	}

	return dc
}

func ProviderPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".rit/providers.json")
}

func CredentialsPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".rit/credentials/")
}

func NewProviderArr(fields Fields) []string {
	var providerArr []string
	for k := range fields {
		if k != AddNew {
			providerArr = append(providerArr, k)
		}
	}
	providerArr = append(providerArr, AddNew)
	return providerArr
}
