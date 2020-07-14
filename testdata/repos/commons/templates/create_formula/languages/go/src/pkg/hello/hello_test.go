package hello

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/gookit/color"
)

func TestHello_Run(t *testing.T) {
	type fields struct {
		Text    string
		List    string
		Boolean string
	}
	tests := []struct {
		name       string
		fields     fields
		wantWriter string
	}{
		{
			name: "Run with success",
			fields: fields{
				Text:    "Hello",
				List:    "World",
				Boolean: "true",
			},
			wantWriter: func() string {
				return fmt.Sprintf("Hello world!\n") +
					color.FgGreen.Render("You receive Hello in text.\n") +
					color.FgRed.Render("You receive World in list.\n") +
					color.FgYellow.Render("You receive true in boolean.\n")
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Hello{
				Text:    tt.fields.Text,
				List:    tt.fields.List,
				Boolean: tt.fields.Boolean,
			}
			writer := &bytes.Buffer{}
			h.Run(writer)
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("Run() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}
