package sliceutil

import "github.com/ZupIT/ritchie-cli/pkg/api"

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
	}
	return false
}
