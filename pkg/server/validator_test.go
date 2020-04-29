package server

import (
	"testing"
)

type finderFoundMock struct{}

func (finderFoundMock) Find() (string, error) {
	return "http://localhost/mocked", nil
}

type finderNotFoundMock struct{}

func (finderNotFoundMock) Find() (string, error) {
	return "", nil
}

func TestValidator(t *testing.T) {

	tests := []struct {
		name string
		in   Finder
		out  error
	}{
		{
			name: "serverURL not found",
			in:   finderNotFoundMock{},
			out:  ErrServerURLNoFound,
		},
		{
			name: "serverURL found",
			in:   finderFoundMock{},
			out:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in
			out := tt.out
			validator := NewValidator(in)

			got := validator.Validate()
			if got != nil && got.Error() != out.Error() {
				t.Errorf("Find(%s) got %v, want %v", in, got, out)
			}
		})
	}
}
