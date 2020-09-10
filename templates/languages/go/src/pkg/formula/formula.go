// This is the formula implementation class.
// Where you will code your methods and manipulate the inputs to perform the specific operation you wish to automate.

package formula

import (
	"fmt"
	"io"

	"github.com/gookit/color"
)

type Formula struct {
	Text     string
	List     string
	Boolean  bool
	Password string
}

func (f Formula) Run(writer io.Writer) {
	var result string

	result += fmt.Sprintf("Hello world!\n")

	result += color.FgGreen.Render(fmt.Sprintf("My name is %s.\n", f.Text))

	if f.Boolean {
		result += color.FgBlue.Render(fmt.Sprintln("I’ve already created formulas using Ritchie."))
	} else {
		result += color.FgRed.Render(fmt.Sprintln("I’m excited in creating new formulas using Ritchie."))
	}

	result += color.FgYellow.Render(fmt.Sprintf("Today, I want to automate %s.\n", f.List))

	result += color.FgCyan.Render(fmt.Sprintf("My secret is %s.\n", f.Password))

	if _, err := fmt.Fprintf(writer, result); err != nil {
		panic(err)
	}
}
