package cmd

import "testing"

func TestNewListCredentialCmd(t *testing.T) {
	cmd := NewListCredentialCmd(singleCredSettingsMock{})
	if cmd == nil {
		t.Errorf("NewListCredentialCmd got %v", cmd)
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}