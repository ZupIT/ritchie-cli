package cmd

import (
	"testing"
)

func TestNewShowContextCmd(t *testing.T) {
	cmd := NewShowContextCmd(ctxFinderMock{}, TutorialFinderMock{})
	if cmd == nil {
		t.Errorf("NewShowContextCmd got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
