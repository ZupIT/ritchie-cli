package cmd

import (
	"testing"
)

func TestNewSingleInitCmd(t *testing.T) {
	cmd := NewSingleInitCmd(inputPasswordMock{}, passphraseManagerMock{}, repoLoaderMock{})
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")

	if cmd == nil {
		t.Errorf("NewSingleInitCmd got %v", cmd)
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}

func TestNewTeamInitCmd(t *testing.T) {
	cmd := NewTeamInitCmd(inputTextMock{}, inputPasswordMock{}, inputURLMock{}, inputFalseMock{}, findSetterServerMock{}, loginManagerMock{}, repoLoaderMock{})
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")

	if cmd == nil {
		t.Errorf("NewTeamInitCmd got %v", cmd)
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
