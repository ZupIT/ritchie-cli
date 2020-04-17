package server

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewValidator(t *testing.T) {
	NewValidator("http://localhost/mocked")
}

func TestValidateSuccess(t *testing.T) {
	url := "http://localhost/mocked"
	v := NewValidator(url)
	got := v.Validate()
	if !reflect.DeepEqual(nil, got) {
		t.Errorf("Set(%s) got %v, want %v", url, got, nil)
	}
}

func TestValidateError(t *testing.T) {
	url := ""
	v := NewValidator(url)
	got := v.Validate()
	if reflect.DeepEqual(nil, got) {
		t.Errorf("Set(%s) got %v, want %v", url, got, fmt.Errorf("No server URL found ! Please set a server URL."))
	}
}