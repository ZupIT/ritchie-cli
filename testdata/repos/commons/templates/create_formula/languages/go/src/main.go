package main

import (
	"hello/pkg/hello"
	"os"
)

func main() {
	input1 := os.Getenv("SAMPLE_TEXT")
	input2 := os.Getenv("SAMPLE_LIST")
	input3 := os.Getenv("SAMPLE_BOOL")

	hello.Hello{
		Text:    input1,
		List:    input2,
		Boolean: input3,
	}.Run(os.Stdout)
}
