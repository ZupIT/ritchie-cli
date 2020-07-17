package set

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

var fileManager = stream.NewFileManager()
var credSettings = credential.NewSettings(fileManager)

func providersPath() string {
	tempDir := os.TempDir()
	path := fmt.Sprintf("%s/providers.json", tempDir)
	return path
}

func TestSingleSettings_ReadCredentials(t *testing.T) {

	credentials, err := credSettings.ReadCredentialsFields("../../../testdata/credentials.json")
	if err != nil {
		t.Errorf("Error on on read credentials function")
	}

	if credentials == nil || len(credentials) <= 0 {
		t.Errorf("Error on on read credentials function, cannot be empty or null")
	}
}

func TestSingleSettings_WriteCredentials(t *testing.T) {
	err := credSettings.WriteCredentialsFields(credential.NewDefaultCredentials(), providersPath())
	defer os.Remove(providersPath())
	if err != nil {
		t.Errorf("Error while write credentials: %s", err)
	}
}

func TestSingleSettings_WriteDefaultCredentials(t *testing.T) {
	err := credSettings.WriteDefaultCredentialsFields(providersPath())
	defer os.Remove(providersPath())
	if err != nil {
		t.Errorf("Error while write credentials: %s", err)
	}
}

func TestNewDefaultCredentials(t *testing.T) {
	defaultCredentials := credential.NewDefaultCredentials()

	if defaultCredentials == nil {
		t.Errorf("Default credentials cannot be nill")
	}

	if len(defaultCredentials) <= 0 {
		t.Errorf("Default credentials cannot be empty")
	}
}

func TestProviderPath(t *testing.T) {
	provider := credential.ProviderPath()
	slicedPath := strings.Split(provider, "/")
	providersJson := slicedPath[len(slicedPath)-1]

	if providersJson != "providers.json" {
		t.Errorf("Providers path must end on providers.json")
	}
}

func TestProvidersArr(t *testing.T) {
	credentials := credential.NewDefaultCredentials()
	providersArray := credential.NewProviderArr(credentials)

	if providersArray[len(providersArray)-1] != credential.AddNew {
		t.Errorf("%q option must be the last one", credential.AddNew)
	}

	if providersArray == nil {
		t.Errorf("Default credentials cannot be nill")
	}

	if len(providersArray) <= 0 {
		t.Errorf("Default credentials cannot be empty")
	}

}
