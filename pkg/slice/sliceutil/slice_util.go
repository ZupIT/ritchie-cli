package sliceutil

import (
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/api"
)

// Contains tells whether a contains s.
func Contains(aa []string, s string) bool {
	for _, v := range aa {
		if s == v {
			return true
		}
	}
	return false
}

// ContainsCmd tells whether a contains c.
func ContainsCmd(aa []api.Command, c api.Command) bool {
	for _, v := range aa {
		if c.Parent == v.Parent && c.Usage == v.Usage {
			return true
		}

		coreCmd := fmt.Sprintf("%s_%s", v.Parent, v.Usage)
		if c.Parent == coreCmd { // Ensures that no core commands will be added
			return true
		}
	}
	return false
}

func Remove(ss []string, r string) []string {
	for i, s := range ss {
		if s == r {
			return append(ss[:i], ss[i+1:]...)
		}
	}
	return ss
}
