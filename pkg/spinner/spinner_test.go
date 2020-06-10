package spinner

import (
	"testing"
	"time"
)

func Test_SpinnerWrapper(t *testing.T) {
	startMsg := "Loading..."
	duration := 1 * time.Millisecond
	s := New(startMsg)

	s.SetSpeed(duration)
	s.SetSleepTime(duration)
	s.Start()
	s.Stop()

	if s.Spinner == nil {
		t.Error("The spinner must not be nil")
	}

	if s.speed != duration {
		t.Errorf("The speed expected %s and got %s", duration, s.speed)
	}

	if s.sleepTime != duration {
		t.Errorf("The sleepTime expected %s and got %s", duration, s.sleepTime)
	}

	if s.Spinner.Title != startMsg {
		t.Errorf("The start message expected %s and got %s", startMsg, s.Spinner.Title)
	}
}
