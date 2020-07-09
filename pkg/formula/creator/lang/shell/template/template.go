package template

const (
	StartFile = "main"

	Main = `#!/bin/sh

. $PWD/{{bin-name}}/{{bin-name}}.sh --source-only

run $SAMPLE_TEXT $SAMPLE_LIST $SAMPLE_BOOL`

	Makefile = `# SH
BINARY_NAME=run.sh
DIST=../bin
FOLDER_SRC={{bin-name}}

build:
	mkdir -p $(DIST)
	cp main.sh $(DIST)/$(BINARY_NAME) && cp -r $(FOLDER_SRC) $(DIST)
	chmod +x $(DIST)/$(BINARY_NAME)`

	Dockerfile = `
FROM alpine:latest

COPY . .

RUN chmod +x set_umask.sh
RUN chmod +x {{bin-name}}.sh
RUN mkdir app

ENTRYPOINT ["./set_umask.sh"]
CMD ["./{{bin-name}}.sh"]`

	File = `#!/bin/sh
run() {
  echo "Hello World! "
  echo "You receive $SAMPLE_TEXT in text. "
  echo "You receive $SAMPLE_LIST in list. " 
  echo "You receive $SAMPLE_BOOL in boolean. "  
}
`
)
