package formula

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/gookit/color"
)

func TestFormula_Run(t *testing.T) {
	type fields struct {
		Text     string
		List     string
		Boolean  bool
		Password string
	}
	tests := []struct {
		name       string
		fields     fields
		wantWriter string
	}{
		{
			name: "Run with TRUE",
			fields: fields{
				Text:     "Dennis",
				List:     "everything",
				Boolean:  true,
				Password: "Ritchie",
			},
			wantWriter: func() string {
				return fmt.Sprintf("Hello world!\n") +
					color.FgGreen.Render(fmt.Sprintf("My name is Dennis.\n")) +
					color.FgBlue.Render(fmt.Sprintln("I've already created formulas using Ritchie.")) +
					color.FgYellow.Render(fmt.Sprintf("Today, I want to automate everything.\n")) +
					color.FgCyan.Render(fmt.Sprintf("My secret is Ritchie.\n"))
			}(),
		},
		{
			name: "Run with FALSE",
			fields: fields{
				Text:     "Dennis",
				List:     "everything",
				Boolean:  false,
				Password: "Ritchie",
			},
			wantWriter: func() string {
				return fmt.Sprintf("Hello world!\n") +
					color.FgGreen.Render(fmt.Sprintf("My name is Dennis.\n")) +
					color.FgBlue.Render(fmt.Sprintln("I'm excited in creating new formulas using Ritchie.")) +
					color.FgYellow.Render(fmt.Sprintf("Today, I want to automate everything.\n")) +
					color.FgCyan.Render(fmt.Sprintf("My secret is Ritchie.\n"))
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Formula{
				Text:     tt.fields.Text,
				List:     tt.fields.List,
				Boolean:  tt.fields.Boolean,
				Password: tt.fields.Password,
			}
			writer := &bytes.Buffer{}
			f.Run(writer)
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("Run() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}
