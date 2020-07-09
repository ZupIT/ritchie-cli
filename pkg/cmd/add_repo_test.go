package cmd

import (
	"net/http"
	"testing"
)

func TestNewAddRepoCmd(t *testing.T) {
	cmd := NewAddRepoCmd(
		&http.Client{},
		repoAdder{},
		inputTextMock{},
		inputPasswordMock{},
		inputURLMock{},
		inputListCredMock{},
		inputTrueMock{},
		inputIntMock{},
	)
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
	t.Skip("Todo test")
	if cmd == nil {
		t.Errorf("NewAddRepoCmd got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}

}
