package credsingle

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

var fileManager = stream.NewFileManager()
var credSettings = NewSingleSettings(fileManager)

func providersPath() string {
	tempDir := os.TempDir()
	path := fmt.Sprintf("%s/providers.json", tempDir)
	return path
}

func TestSingleSettings_ReadCredentials(t *testing.T) {

	credentials, err := credSettings.ReadCredentials("../../../testdata/credentials.json")
	if err != nil {
		t.Errorf("Error on on read credentials function")
	}

	if credentials == nil || len(credentials) <= 0 {
		t.Errorf("Error on on read credentials function, cannot be empty or null")
	}
}

func TestSingleSettings_WriteCredentials(t *testing.T) {
	err := credSettings.WriteCredentials(NewDefaultCredentials(), providersPath())
	defer os.Remove(providersPath())
	if err != nil {
		t.Errorf("Error while write credentials: %s", err)
	}
}

func TestSingleSettings_WriteDefaultCredentials(t *testing.T) {
	credentials := credential.Fields{
		"customField":     []credential.Field{},
	}
	fieldsData, err := json.Marshal(credentials)
	if err != nil {
		t.Errorf("Error while writing existing credentials: %s", err)
	}

	// Write an initial credential file
	err = ioutil.WriteFile(providersPath(), fieldsData, os.ModePerm)
	defer os.Remove(providersPath())
	if err != nil {
		t.Errorf("Error while writing existing credentials: %s", err)
	}

	// Call the method
	err = credSettings.WriteDefaultCredentials(providersPath())
	if err != nil {
		t.Errorf("Error while writing existing credentials: %s", err)
	}

	// Reopen file and check if previous config was not lost
	file, _ := ioutil.ReadFile(providersPath())
	var fields credential.Fields
	err = json.Unmarshal(file, &fields)
	if err != nil {
		t.Errorf("Error while writing existing credentials: %s", err)
	}
	if len(fields) != len(NewDefaultCredentials()) + 1 {
		t.Errorf("Writing existing credentials did not succeed in adding a field")
	}
	if fields["customField"] == nil {
		t.Errorf("Writing existing credentials did not save custom field")
	}
}

func TestSingleSettings_WriteDefaultCredentialsOnExistingFile(t *testing.T) {
	err := credSettings.WriteDefaultCredentials(providersPath())
	defer os.Remove(providersPath())
	if err != nil {
		t.Errorf("Error while write credentials: %s", err)
	}
}

func TestNewDefaultCredentials(t *testing.T) {
	defaultCredentials := NewDefaultCredentials()

	if defaultCredentials == nil {
		t.Errorf("Default credentials cannot be nill")
	}

	if len(defaultCredentials) <= 0 {
		t.Errorf("Default credentials cannot be empty")
	}
}

func TestProviderPath(t *testing.T) {
	provider := ProviderPath()
	slicedPath := strings.Split(provider, "/")
	providersJson := slicedPath[len(slicedPath)-1]

	if providersJson != "providers.json" {
		t.Errorf("Providers path must end on providers.json")
	}
}

func TestProvidersArr(t *testing.T) {
	credentials := NewDefaultCredentials()
	providersArray := NewProviderArr(credentials)

	if providersArray[len(providersArray)-1] != AddNew {
		t.Errorf("%q option must be the last one", AddNew)
	}

	if providersArray == nil {
		t.Errorf("Default credentials cannot be nill")
	}

	if len(providersArray) <= 0 {
		t.Errorf("Default credentials cannot be empty")
	}

}
