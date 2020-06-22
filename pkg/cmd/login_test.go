package cmd

import (
	"testing"
)

func TestNewLoginCmd(t *testing.T) {
	cmd := NewLoginCmd(inputTextMock{}, inputPasswordMock{}, loginManagerMock{}, repoLoaderMock{}, findSetterServerMock{})
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
	if cmd == nil {
		t.Errorf("NewLoginCmd got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
