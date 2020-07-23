package modifier

import "github.com/ZupIT/ritchie-cli/pkg/formula"

type Modifier interface {
	modify(b []byte) []byte
}

func NewModifiers(create formula.Create) []Modifier {
	return []Modifier{
		FormulaCmd{cf: create},
		FormulaTags{cf: create},
	}
}

func Modify(b []byte, modifiers []Modifier) []byte {
	result := b
	for _, m := range modifiers {
		result = m.modify(result)
	}
	return result
}
