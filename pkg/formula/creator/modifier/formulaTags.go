package modifier

import (
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

type FormulaTags struct {
	cf formula.Create
}

func (f FormulaTags) modify(b []byte) []byte {

	content := string(b)
	tags := toTags(f.cf.FormulaCmdName())
	content = strings.ReplaceAll(content, "#rit-replace{formulaTags}", tags)
	return []byte(content)

}

func toTags(fCmdName string) string {
	return strings.ReplaceAll(fCmdName, " ", `", "`)
}
