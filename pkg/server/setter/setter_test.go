package setter

import (
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/api"
)

func TestNewSetter(t *testing.T) {
	NewSetter(api.RitchieHomeDir())
}

func TestSet(t *testing.T) {
	s := NewSetter(api.RitchieHomeDir())
	err := s.Set("http://localhost/mocked"); if err != nil { return }
}