package cmd

import (
	"testing"
)

func TestNewSingleSetCredentialCmd(t *testing.T) {
	cmd := NewSetCredentialCmd(credSetterMock{}, singleCredSettingsMock{}, inputSecretMock{}, inputFalseMock{}, inputListCredMock{}, inputPasswordMock{})

	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
	if cmd == nil {
		t.Errorf("NewSetCredentialCmd got %v", cmd)
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}

// Todo fix TestNewTeamSetCredentialCmd
// func TestNewTeamSetCredentialCmd(t *testing.T) {
// 	cmd := NewTeamSetCredentialCmd(credSetterMock{}, credSettingsMock{}, inputSecretMock{}, inputFalseMock{}, inputListCredMock{}, inputPasswordMock{}, InputMultilineMock{})
// 	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
// 	if cmd == nil {
// 		t.Errorf("NewTeamSetCredentialCmd got %v", cmd)
// 	}
//
// 	if err := cmd.Execute(); err != nil {
// 		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
// 	}
// }
