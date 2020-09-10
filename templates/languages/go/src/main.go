// This is the main class.
// Where you will extract the inputs asked on the config.json file and call the formula's method(s).

package main

import (
	"formula/pkg/formula"
	"os"
	"strconv"
)

func main() {
	input1 := os.Getenv("INPUT_TEXT")
	input2 := os.Getenv("INPUT_LIST")
	input3, _ := strconv.ParseBool(os.Getenv("INPUT_BOOLEAN"))
	input4 := os.Getenv("INPUT_PASSWORD")

	formula.Formula{
		Text:     input1,
		List:     input2,
		Boolean:  input3,
		Password: input4,
	}.Run(os.Stdout)
}
