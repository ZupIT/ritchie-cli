package template

const (
	StartFile = "main"

	Main = `mod {{bin-name}};

use std::env;
	
fn main() {
	let sample_text;
	let sample_list;
	let sample_bool;

	match env::var("SAMPLE_TEXT") {
		Ok(val) => sample_text = val,
		Err(_e) => sample_text = "none".to_string(),
	}

	match env::var("SAMPLE_LIST") {
		Ok(val) => sample_list = val,
		Err(_e) => sample_list = "none".to_string(),
	}

	match env::var("SAMPLE_BOOL") {
		Ok(val) => sample_bool = val,
		Err(_e) => sample_bool = "none".to_string(),
	}

	{{bin-name}}::run(sample_text, sample_list, sample_bool);
}`

	Run = `#!/bin/sh

cargo install --path .

cargo run`

	File = `pub fn run(sample_text: String, sample_list: String, sample_bool: String) {
	println!("Hello World!");
	println!("You receive {} in text.", sample_text);
  println!("You receive {} in list.", sample_list);
  println!("You receive {} in boolean.", sample_bool);
}`

	CargoToml = `[package]
name = "rust"
version = "0.1.0"
edition = "2018"
	
[dependencies]`

	Dockerfile = `
FROM rust:latest
WORKDIR /app

COPY . .

RUN chmod +x set_umask.sh

RUN cargo install --path .

ENTRYPOINT ["/set_umask.sh"]
CMD ["app"]
`

	Makefile = `# Make Run Rust
BINARY_NAME_UNIX={{bin-name}}.sh
BINARY_NAME_WINDOWS={{bin-name}}.bat
DIST=../dist
DIST_DIR=$(DIST)/commons/bin

build:
	mkdir -p $(DIST_DIR)
	cp run_template $(BINARY_NAME_UNIX) && chmod +x $(BINARY_NAME_UNIX)
	sed '1d' run_template > $(BINARY_NAME_WINDOWS) && chmod +x $(BINARY_NAME_WINDOWS)

	cp -r . $(DIST_DIR)

	#Clean files
	rm $(BINARY_NAME_UNIX)`

	WindowsBuild = `:: Rust parameters
echo off
SETLOCAL
SET BINARY_NAME_UNIX={{bin-name}}.sh
SET BINARY_NAME_WINDOWS={{bin-name}}.bat
SET DIST=..\dist
SET DIST_DIR=%DIST%\commons\bin
:build
    mkdir %DIST_DIR%
	more +1 run_template > %DIST_DIR%\%BINARY_NAME_WINDOWS%
    copy run_template %DIST_DIR%\%BINARY_NAME_UNIX%
    xcopy . %DIST_DIR% /E /H /C /I
    GOTO DONE
:DONE`
)
