package credsingle

import (
	"strings"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

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
	fileManager := stream.NewFileManager()
	credSettings := NewSingleSettings(fileManager)
	credentials, err := credSettings.ReadCredentials("../../../testdata/credentials.json")
	if err != nil {
		t.Errorf("Error on on read credentials function")
	}

	if credentials == nil || len(credentials) <= 0 {
		t.Errorf("Error on on read credentials function, cannot be empty or null")
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

