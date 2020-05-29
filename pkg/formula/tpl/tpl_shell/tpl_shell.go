package tpl_shell

const (
	Main = `#!/bin/sh

. ./{{bin-name}}/{{bin-name}}.sh --source-only

run $SAMPLE_TEXT $SAMPLE_LIST $SAMPLE_BOOL`

	Makefile = `# SH
BINARY_NAME={{bin-name}}.sh
DIST=../dist
DIST_DIR=$(DIST)/commons/bin
build:
	mkdir -p $(DIST_DIR)
	cp main.sh $(DIST_DIR)/$(BINARY_NAME) && cp -r {{bin-name}} $(DIST_DIR) && cp Dockerfile $(DIST_DIR)
	chmod +x $(DIST_DIR)/$(BINARY_NAME)`

	Dockerfile = `
FROM alpine:latest

WORKDIR /app

COPY . .

RUN chmod +x main.sh
RUN chmod +x main
RUN chmod +x set_umask.sh

ENTRYPOINT ["/set_umask.sh"]
CMD /app/main.sh`

	File = `#!/bin/sh
run() {
  echo "Hello World! "
  echo "You receive $SAMPLE_TEXT in text. "
  echo "You receive $SAMPLE_LIST in list. " 
  echo "You receive $SAMPLE_BOOL in boolean. "  
}
`

	Umask = `#!/bin/sh
umask 0011
$1
`
)
