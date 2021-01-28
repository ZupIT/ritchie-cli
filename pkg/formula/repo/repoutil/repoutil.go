package repoutil

import (
	"fmt"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

const localPrefix = "local-%s"

func LocalName(repoName string) formula.RepoName {
	repoName = strings.ToLower(fmt.Sprintf(localPrefix, repoName))
	return formula.RepoName(repoName)
}
