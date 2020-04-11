package cmd

import (
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/autocomplete"
)

func TestNewAutocompleteCmd(t *testing.T) {
	cmd := NewAutocompleteCmd()
	if cmd == nil {
		t.Errorf("NewAutocompleteCmd got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}

type autocompleteGenMock struct{}

func (autocompleteGenMock) Generate(s autocomplete.ShellName) (string, error) {
	return "autocomplete", nil
}

func TestNewAutocompleteZsh(t *testing.T) {
	mock := autocompleteGenMock{}
	cmd := NewAutocompleteZsh(mock)
	if cmd == nil {
		t.Errorf("NewAutocompleteZsh got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}

func TestNewAutocompleteBash(t *testing.T) {
	mock := autocompleteGenMock{}
	cmd := NewAutocompleteBash(mock)
	if cmd == nil {
		t.Errorf("NewAutocompleteBash got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
