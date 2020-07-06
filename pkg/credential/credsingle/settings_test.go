package credsingle

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

var fileManager = stream.NewFileManager()
var credSettings = NewSingleSettings(fileManager)

func TestNewDefaultCredentials(t *testing.T) {
	defaultCredentials := NewDefaultCredentials()

	if defaultCredentials == nil {
		t.Errorf("Default credentials cannot be nill")
	}

	if len(defaultCredentials) <= 0 {
		t.Errorf("Default credentials cannot be empty")
	}
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
	home, _ := os.UserHomeDir()
	err := credSettings.WriteCredentials(NewDefaultCredentials(), home)
	if err != nil {
		t.Errorf("Error while write credentials: %s", err)
	}

	os.Remove(fmt.Sprintf("%s/providers.json", home))
}

func TestProviderPath(t *testing.T) {
	provider := ProviderPath()
	slicedPath := strings.Split(provider, "/")
	providersJson := slicedPath[len(slicedPath)-1]

	if providersJson != "providers.json" {
		t.Errorf("Providers path must end on providers.json")
	}
}
