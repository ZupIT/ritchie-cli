package cmd

import (
	"testing"
)

func TestNewAddRepoCmd(t *testing.T) {
	cmd := NewAddRepoCmd(repoAdder{}, inputTextMock{}, inputURLMock{}, inputIntMock{}, inputTrueMock{})
	if cmd == nil {
		t.Errorf("NewAddRepoCmd got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}

}
