package cmd

import (
	"errors"
	"testing"
)

func TestNewSingleSetCredentialCmd(t *testing.T) {
	cmd := NewSingleSetCredentialCmd(credSetterMock{}, inputSecretMock{}, inputFalseMock{}, inputListCredMock{}, inputPasswordMock{})
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
	if cmd == nil {
		t.Errorf("NewSingleSetCredentialCmd got %v", cmd)
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}

func TestNewTeamSetCredentialCmd(t *testing.T) {
	cmd := NewTeamSetCredentialCmd(credSetterMock{}, credSettingsMock{}, inputSecretMock{}, inputFalseMock{}, inputListCredMock{}, inputPasswordMock{})
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
	if cmd == nil {
		t.Errorf("NewTeamSetCredentialCmd got %v", cmd)
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}

func TestNewSingleSetCredentialCmdWithEntryArchive(t *testing.T) {
	errPath := errors.New("some error of path")
	cmd := NewSingleSetCredentialCmd(
		credSetterMock{},
		inputTextCustomMock{
			text: func(name string, required bool) (string, error) {
				if name == MsgTypeEntryPath {
					return "/test", errPath
				}
				return "some_path", nil
			},
		},
		inputFalseMock{},
		inputListCustomMock{
			list: func(name string, list []string) (string, error) {
				if name == MsgTypeEntry {
					return EntriesTypeCredentialFile, nil
				}
				return "some_input", nil
			},
		},
		inputPasswordMock{},
	)
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
	if cmd == nil {
		t.Errorf("NewTeamSetCredentialCmd got %v", cmd)
	}

	if err := cmd.Execute(); err != errPath {
		t.Errorf("%s = '%v', want '%v'", cmd.Use, err, errPath)
	}
}
