package setter

import (
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
)

const serverFilePattern = "%s/server"

type Setter struct {
	serverFile string
}

func NewSetter(ritchieHomeDir string) Setter {
	return Setter{
		serverFile: fmt.Sprintf(serverFilePattern, ritchieHomeDir),
	}
}

func (s Setter) Set(url string) error {
	if err := fileutil.WriteFile(s.serverFile, []byte(url)); err != nil {
		return err
	}
	return nil
}
