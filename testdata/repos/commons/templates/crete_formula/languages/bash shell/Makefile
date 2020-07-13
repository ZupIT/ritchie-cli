# SH
BINARY_NAME=run.sh
DIST=bin

build: bash-build

bash-build:
	mkdir -p $(DIST)
	cp -r src/* $(DIST)
	mv $(DIST)/main.sh $(DIST)/$(BINARY_NAME)
	chmod +x $(DIST)/$(BINARY_NAME)
