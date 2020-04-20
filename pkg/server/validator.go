package server

import (
	"fmt"
)

type ValidatorManager struct {
	serverFinder Finder
}

func NewValidator(serverFinder Finder) Validator{
	return ValidatorManager{
		serverFinder: serverFinder,
	}
}

func (v ValidatorManager) Validate() error {
	serverUrl, err := v.serverFinder.Find()
	if err != nil {
		return err
	}

	if serverUrl == "" {
		return fmt.Errorf("No server URL found ! Please set a server URL.")
	}

	return nil
}
