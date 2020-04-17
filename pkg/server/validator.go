package server

import (
	"fmt"
)

type ValidatorManager struct {
	serverUrl string
}

func NewValidator(serverUrl string) Validator{
	return ValidatorManager{
		serverUrl: serverUrl,
	}
}

func (v ValidatorManager) Validate() error {
	if v.serverUrl == "" {
		return fmt.Errorf("No server URL found ! Please set a server URL.")
	}
	return nil
}
