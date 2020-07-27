package modifier

import (
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

type FormulaCmd struct {
	cf formula.Create
}

func (f FormulaCmd) modify(b []byte) []byte {

	content := string(b)
	content = strings.ReplaceAll(content, "#rit-replace{formulaCmd}", f.cf.FormulaCmd)
	return []byte(content)

}
