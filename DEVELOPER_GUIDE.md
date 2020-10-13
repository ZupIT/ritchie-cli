# Developer guide

Digging in the project for the first time might make it difficult to understand where to start. Hopefully, this
guide will provide you with a better understanding of the project structure. 

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

#### Internal

Ritchie collects anonymous usage metrics with user consent, so we can understand user behavior and always
keep improving the cli. Commands are sent to a database using protocol buffers via grpc connection. 

#### Packaging

Contains code and instructions to package Ritchie for different supported OS distributions. Here we have scripts
for the Windows installer, rpm and debian distributions, and other installation scripts.

#### Pkg

This module is the core of the project, any command or functionality is developed here and added to **Cobra**
via `main`. 

#### Testdata

Contains multiple dummy files used for tests.