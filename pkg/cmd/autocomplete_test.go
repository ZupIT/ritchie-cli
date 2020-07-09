package cmd

import (
	"testing"
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

func TestNewAutocompleteFish(t *testing.T) {
	mock := autocompleteGenMock{}
	cmd := NewAutocompleteFish(mock)
	if cmd == nil {
		t.Errorf("NewAutocompleteFish got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}

func TestNewAutocompletePowerShell(t *testing.T) {
	mock := autocompleteGenMock{}
	cmd := NewAutocompletePowerShell(mock)
	if cmd == nil {
		t.Errorf("NewAutocompletePowerShell got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
