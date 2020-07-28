package credential

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

var fileManager = stream.NewFileManager()
var dirManager = stream.NewDirManager(fileManager)
var homeDir, _ = os.UserHomeDir()
var credSettings = NewSettings(fileManager, dirManager, homeDir)

func providersPath() string {
	tempDir := os.TempDir()
	path := fmt.Sprintf("%s/providers.json", tempDir)
	return path
}

func TestSettings_ReadCredentialsFields(t *testing.T) {
	credentials, err := credSettings.ReadCredentialsFields("../../testdata/credentials.json")
	if err != nil {
		t.Errorf("Error reading credentials fields")
	}

	if credentials == nil || len(credentials) <= 0 {
		t.Errorf("Error reading credentials fields, cannot be empty or null")
	}
}

func TestSettings_ReadCredentialsValue(t *testing.T) {
	credentials, err := credSettings.ReadCredentialsValue("../../testdata/.rit/credentials/")
	if err != nil {
		t.Errorf("Error reading credentials: %s", err)
	}

	if credentials == nil || len(credentials) <= 0 {
		t.Errorf("Error reading credentials, cannot be empty or null")
	}
}

func TestSettings_WriteCredentialsFields(t *testing.T) {
	defer os.Remove(providersPath())
	var tests = []struct {
		name    string
		path    string
		fields  Fields
		wantErr bool
	}{
		{
			name:    "Run with success",
			path:    providersPath(),
			fields:  NewDefaultCredentials(),
			wantErr: false,
		},
		{
			name:    "Error with invalid path",
			path:    "",
			fields:  NewDefaultCredentials(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := credSettings.WriteCredentialsFields(tt.fields, tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write credentials fields error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

func TestSettings_WriteDefaultCredentialsFields(t *testing.T) {
	err := credSettings.WriteDefaultCredentialsFields(providersPath())
	defer os.Remove(providersPath())
	if err != nil {
		t.Errorf("Error writing credentials: %s", err)
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
	provider := credSettings.ProviderPath()
	slicedPath := strings.Split(provider, "/")
	providersJson := slicedPath[len(slicedPath)-1]

	if providersJson != "providers.json" {
		t.Errorf("Providers path must end on providers.json")
	}
}

func TestCredentialsPath(t *testing.T){
	credentials := credSettings.CredentialsPath()
	slicedPath := strings.Split(credentials, "/")
	fmt.Println(slicedPath)
	providersDir := slicedPath[len(slicedPath)-1]

	if providersDir != "credentials"{
		t.Errorf("Providers path must end on credentials dir")
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
