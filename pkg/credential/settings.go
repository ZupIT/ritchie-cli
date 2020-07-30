package credential

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

const AddNew = "Add a new"

type Settings struct {
	file    stream.FileWriteReadExistLister
	dir     stream.DirLister
	HomeDir string
}

func NewSettings(file stream.FileWriteReadExistLister, dir stream.DirLister, homeDir string) Settings {
	return Settings{
		file:    file,
		dir:     dir,
		HomeDir: homeDir,
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
			cred.Credential = formatCredential(string(cBytes))
			cred.Provider = detail.Service
			cred.Context = c
			creds = append(creds, cred)
			detail = Detail{}
		}
	}
	return creds, nil
}

func formatCredential(credential string) string {
	credArr := strings.Split(credential, "credential")
	credArr = strings.Split(credArr[1], "service")

	credValue := strings.TrimPrefix(credArr[0], "\":")
	credValue = strings.TrimSuffix(credValue, ",\"")

	splitedCredential := strings.Split(credValue, "\"")
	for i, c := range splitedCredential {
		if c == ":" {
			splitedCredential[i+1] = formatCredValue(splitedCredential[i+1])
		}
	}

	return strings.Join(splitedCredential, "\"")
}

func formatCredValue(credential string) string {
	if credLen := len(credential); credLen > 12 {
		var resumedCredential []rune
		for i, r := range credential {
			if i >= 4 {
				r = '*'
			}
			resumedCredential = append(resumedCredential, r)
			if i > 10 {
				break
			}
		}
		return string(resumedCredential) + "..."
	} else {
		var hiddenCredential []rune
		mustHideIndex := credLen / 3
		for i, r := range credential {
			if i > mustHideIndex {
				r = '*'
			}
			hiddenCredential = append(hiddenCredential, r)
		}
		return string(hiddenCredential)
	}
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
		Type: "plain text",
	}

	token := Field{
		Name: "token",
		Type: "secret",
	}

	accessKey := Field{
		Name: "accesskey",
		Type: "plain text",
	}

	secretAccessKey := Field{
		Name: "secretaccesskey",
		Type: "secret",
	}

	base64config := Field{
		Name: "base64config",
		Type: "plain text",
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

func (s Settings) ProviderPath() string {
	return filepath.Join(s.HomeDir, ".rit/providers.json")
}

func (s Settings) CredentialsPath() string {
	return filepath.Join(s.HomeDir, ".rit/credentials/")
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
