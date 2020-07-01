package credsingle

import (
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

type Settings struct {
	file stream.FileWriteReadExister
}

func NewSettings(file stream.FileWriteReadExister) Settings {
	return Settings{file: file}
}



// func write json
// homeDir, _ := os.UserHomeDir()
// providerDir := fmt.Sprintf("%s/.rit/repo/providers.json", homeDir)

// credentialData, _ := json.Marshal()
// _ = fileutil.WriteFile(providerDir, credentialData)
