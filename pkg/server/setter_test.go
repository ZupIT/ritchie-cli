package server

import (
	"errors"
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/stream"
	"github.com/ZupIT/ritchie-cli/pkg/validator"
)

func TestSet(t *testing.T) {

	s := NewSetter(os.TempDir(), stream.NewFileManager())

	tests := []struct {
		name string
		in   string
		out  error
	}{
		{
			name: "empty serverURL",
			in:   "",
			out:  validator.ErrInvalidServerURL,
		},
		{
			name: "existing serverURL",
			in:   "http://localhost/mocked",
			out:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.in
			out := tt.out
			got := s.Set(in)
			if got != nil && errors.Unwrap(got).Error() != out.Error() {
				t.Errorf("Set(%s) got %v, want %v", in, got, out)
			}
		})
	}
}
