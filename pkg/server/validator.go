package server

import (
	"errors"
)

// ErrServerURLNoFound when the serverURL is not found
var ErrServerURLNoFound = errors.New("not found server URL, please set one")

type ValidatorManager struct {
	serverFinder Finder
}

func NewValidator(serverFinder Finder) Validator {
	return ValidatorManager{
		serverFinder: serverFinder,
	}
}

func (v ValidatorManager) Validate() error {
	url, err := v.serverFinder.Find()
	if err != nil {
		return err
	}

	if url == "" {
		return ErrServerURLNoFound
	}

	return nil
}
