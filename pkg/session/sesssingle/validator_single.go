package sesssingle

import "github.com/ZupIT/ritchie-cli/pkg/session"

type Validator struct {
	manager session.Manager
}

func NewValidator(m session.Manager) Validator {
	return Validator{m}
}

func (s Validator) Validate() error {
	_, err := s.manager.Current()
	return err
}
