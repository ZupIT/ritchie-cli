package server

import (
	"os"
	"reflect"
	"testing"
)

func TestNewSetter(t *testing.T) {
	NewSetter(os.TempDir())
}

func TestSetSuccess(t *testing.T) {
	s := NewSetter(os.TempDir())
	url := "http://localhost/mocked"
	got := s.Set(url)
	if !reflect.DeepEqual(nil, got) {
		t.Errorf("Set(%s) got %v, want %v", url, got, nil)
	}
}

func TestSetPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	s := NewSetter(os.TempDir())
	url := ""
	s.Set(url)
}