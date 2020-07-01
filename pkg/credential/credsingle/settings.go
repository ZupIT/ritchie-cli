package credsingle

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

var file stream.FileWriteReadExister

func ReadCredentialsJson() credential.Fields{
	var fields credential.Fields

	cBytes, _ := file.Read(providerPath())
	_ = json.Unmarshal(cBytes, fields)

	return fields
}

func WriteCredentialsJson(fields credential.Fields) error {
	fieldsData, _  := json.Marshal(fields)

	err := file.Write(providerPath(), fieldsData); if err != nil {
		return err
	}

	return nil
}

func providerPath() string{
	homeDir, _ := os.UserHomeDir()
	providerDir := fmt.Sprintf("%s/.rit/repo/providers.json", homeDir)
	return providerDir
}
// func write json
