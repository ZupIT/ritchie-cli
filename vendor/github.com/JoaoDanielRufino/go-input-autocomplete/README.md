# go-input-autocomplete

A useful input that can autocomplete users path to directories or files when tab key is pressed. The purpose is to be similar to bash/cmd native autocompletion.

## Installation

```bash
go get github.com/JoaoDanielRufino/go-input-autocomplete
```

## Usage

```go
package main

import (
	"fmt"
	input_autocomplete "github.com/JoaoDanielRufino/go-input-autocomplete"
)

func main() {
	path, err := input_autocomplete.Read("Path: ")

	if err != nil {
		panic(err)
	}

	fmt.Println(path)
}
```

## How it works

![gif](doc/go-input-autocomplete.gif)
