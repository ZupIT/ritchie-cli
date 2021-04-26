# Developer guide

Digging in the project for the first time might make it difficult to understand where to start. Hopefully, this
guide will provide you with a better understanding of the project structure. 

## Running the project

There are some ways to setup the development environment. You can run directly from your favorite
IDE or overwrite Ritchie's binary on your machine. The most popular IDEs are [GoLand](https://www.jetbrains.com/pt-br/go/)
from JetBrains and [VSCode](https://code.visualstudio.com/) from Windows. 

#### Running on GoLand

* On the top menu go to `Run > Edit Configurations`. Click on the `+` button on the top left and select to create a 
new `Go Build`.

* On **Files**, select the `cmd/main.go` file as the entry point of your program.

* As for **Working Directory**, use the project root

* Under **Program Arguments**, write the command you want to run without the `rit`. For instance, if you want to 
debug `rit create formula`, just add `create formula` on this field. 

<img class="special-img-class" src="/docs/img/go-land-setup.png"  alt="Edit Configurations screen of GoLand"/>

#### Overwriting your current binary

If you want to basically use your local code as the running Ritchie distribution on your machine, just run this 
simple command from your project root directory in the terminal:

* Linux: `make build-linux && sudo mv ./dist/linux/rit /usr/local/bin`
* Mac: `make build-mac && sudo mv ./dist/darwin/rit /usr/local/bin`
* Windows: `make build-windows && runas \user:administrator move .\dist\windows\rit.exe "\Program Files (x86)\rit"`

## The main libraries

There are two main libraries that are most used throughout the project. These are the ones that aid us directly
in building a solid cli, and probably the ones you will have most contact with.

#### [Cobra](https://github.com/spf13/cobra)

Cobra is a powerful cli building tool. Many modern projects use it to build their clis. We can map out all commands
and subcommands in Ritchie and Cobra will take care of executing them with the right methods and flags, provide
helper descriptions for users, and even more features. Check its github repo to learn more about it. 

#### [Survey](https://github.com/AlecAivazis/survey)

One of the key features of Ritchie is that it lets you run formulas in an interactive manner. To help us achieve
a better user experience, we use survey. Survey is a library of cli prompts in various formats. It supports text
inputs, secret inputs, list selection, among other things.

## Project structure

For this section, we discuss the main folders in the project and their purpose.

#### Cmd

This is where the `main` file resides. Ritchie comes natively with a set of **core** commands to manipulate the cli.
These commands range from managing workspaces and credentials, to building and creating formulas. We make heavy 
use of dependency injection to make these commands testable. Therefore, most of the components are initialized here
and then added to `Cobra`. If you have any new command to add, you will be definitely editing this file. Keep in mind
the organisation of the commands!

#### Functional

Currently, we have two main ways of running formulas. Formulas can be run interactively or via stdin. This folder
contains the files to run formulas on the different supported OSes and with different input formats.

#### Packaging

Contains code and instructions to package Ritchie for different supported OS distributions. Here we have scripts
for the Windows installer, rpm and debian distributions, and other installation scripts.

#### Pkg

This module is the core of the project, any command or functionality is generally developed here and added to **Cobra**
via `main`. Some of the features developed are listed below:
* api: contains main constants such as the user home directory and core commands list.
* autocomplete: contains inline scripts and logic for the `rit completion` commands.
* cmd: contains most of core commands implementation. You can find commands such as listing, creating, 
and deleting resources here. For simple implementations they are enough, but they might call dedicated modules
to perform certain actions.
* credential: contains logic related to credential manipulation, such as the files they are saved and how to manage
them.
* env: helper to resolve credential from input runners.
* formula: contains all formula manipulation implementations, such as building, resolving, creating, 
and running formulas
* git: module to add, remove, or manage other formula repos.
* http: header definitions
* metric: sends collected anonymous metrics via http requests.
* prompt: Ritchie's adaptations on the `survey` module for user input.
* rcontext: manages user contexts. Handles files modifications and context activation so users can run formulas
using different sets of credentials (i.e.: development, staging, production)
* rtutorial: basic implementation of tutorial texts. Tutorial are helper texts that can be added to each command
to provide the user with more context on that action.
* stdin: JSON stdin parser
* upgrade: manages and performs upgrades on Ritchie
* version: manages and resolves Ritchie's versioning

#### Testdata

Contains multiple dummy files used for tests.

## Testing your code

Ritchie relies heavily on dependency injection to test its code. If you are creating a new command, you should
pass on any required dependencies from the `main.go` initialization file. Usually you pass on inputs to prevent
the test from blocking the execution flow while awaiting for an user input. For instance, suppose you have a list
that offers the options `a, b, c`. Then, you simply pass a mock that returns `a` straight away in place of the listing
struct to perform your test.

#### Manipulations on your file system

When testing, think about the intended behavior of your command. Does it create or modify a file? Does it manipulate
a directory? You can make use of the `os.TempDir()` command to setup your sandbox file system and manipulate files
and folders there. Let's suppose your implementation edits the credential config file. You do not want to meddle
with your own config file every time you run a test. So simply create one in a temp directory and pass its path
to the command *as if it was ritchie's home directory.

Do not forget to remove your temp directory to not flood your file system with trash on each run! You can use `defer`
to ensure a piece of code will always be executed after the method returns regardless of where the execution ended. 

`defer os.RemoveAll(filePath)`

#### Test structure

Test files are suffixed with `_test`. Test functions have the format `TestCamelCase` which will be printed when 
you execute them. If your test has multiple outcome branches, you can create a struct array to iterate over such cases.
The struct usually has the format:
```go
    tests := []struct {
		name    string   // Name of the test to display
		fields  fields   // Contains all necessary struct mocks or values to perform the test
		output  bool     // Expected result to compare to
	}
```

A typical test has the following structure
```go
    // Mocks initialization here
    
	var tests = []struct {
		name            string
		fields          fields
		output          string
	}{
		{
			name:            "Some name",
			fields:          fields{ ... },
			output:          "Some output",
		},
		...,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Your execution code here that compares with the output
		})
	}
```

> Note: we are currently reviewing how our tests are written, we are trying to make a more assertive test framework
>on Ritchie cli, such as simulating the user input via stdin instead of mocking it, so these practices might change.


## Translation

- All translation files are inside the `/resources/i18n` folder.
- Translation files follow a convention, must start with your locale, and have the extension `.toml`. For example: `en.toml`
- Any phrase or word that does not exist in `en.toml` must be added to it to avoid mistakes.
- The convention for translation IDs is to start with the name of the specific package or command followed by a clear identifier for your word or phrase. For example: `[init.welcome]`

### See some examples below:

#### English file (`en.toml`)

```toml
[init.hello.world]
other = "Hello World!"
```  

#### Portuguese file (`pt_BR.toml`)

```toml
[init.hello.world]
other = "Olá Mundo!"
```  

#### Other locale example (`ge_GE.toml`)

```toml
[init.hello.world]
other = "Hallo Welt!"
```

> Note: you must run the `make generate-translation` command to update the `translations.go` file whenever you create 
> or update a translation. 

### How to use translation in your Golang code?

You will use the `i18n` package located inside `internal/pkg/i18n`. See an example below:

```go
import "github.com/ZupIT/ritchie-cli/internal/pkg/i18n"


var helloWorldMsg = i18n.T("hello.world") // [hello.world] is the message ID.

```

> :warning: Attention: For new languages you must add the full name of the language and its locale on a map called `i18n.Langs`.

Example:

``` go
package i18n

var Langs = map[string]string{
	"English":    "en",
	"Portuguese": "pt_BR",
	"German":     "ge_GE", // Must be the same name as the translation file
}

```