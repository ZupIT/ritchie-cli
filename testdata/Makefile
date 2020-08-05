#Makefiles
TERRAFORM=aws/terraform
NAVIGATE_HANDBOOK=github/navigate-handbook
SEARCH_HANDBOOK=github/search-handbook
SC_COFFEE_GO=scaffold/coffee-go
SC_COFFEE_JAVA=scaffold/coffee-java
SC_COFFEE_SHELL=scaffold/coffee-shell
SC_COFFEE_PHP=scaffold/coffee-php
SC_COFFEE_PYTHON=scaffold/coffee-python
SC_COFFEE_NODE=scaffold/coffee-node
SC_SPRING_STARTER=scaffold/spring-starter
KAFKA=kafka
DOCKER=docker/compose
KUBERNETES=kubernetes/core
FAST_MERGE=github/fast-merge

FORMULAS=$(TERRAFORM) $(SC_COFFEE_GO) $(SC_COFFEE_JAVA) $(SC_COFFEE_SHELL) $(SC_COFFEE_PHP) $(SC_COFFEE_PYTHON) $(SC_COFFEE_NODE) $(SC_SPRING_STARTER) $(KAFKA) $(DOCKER) $(NAVIGATE_HANDBOOK) $(SEARCH_HANDBOOK) $(KUBERNETES) $(FAST_MERGE)

PWD_INITIAL=$(shell pwd)

FORM_TO_UPPER  = $(shell echo $(form) | tr  '[:lower:]' '[:upper:]')
FORM = $($(FORM_TO_UPPER))

push-s3:
	echo $(RITCHIE_AWS_BUCKET)
	echo "Init pwd: $(PWD_INITIAL)"
	for formula in $(FORMULAS); do cd $$formula/src && make build && cd $(PWD_INITIAL) || exit; done
	./copy-bin-configs.sh "$(FORMULAS)"
	aws s3 cp . s3://$(RITCHIE_AWS_BUCKET)/ --exclude "*" --include "formulas/*" --recursive
	aws s3 cp . s3://$(RITCHIE_AWS_BUCKET)/ --exclude "*" --include "tree/tree.json" --recursive
	rm -rf formulas

bin:
	echo "Init pwd: $(PWD_INITIAL)"
	echo "Formulas bin: $(FORMULAS)"
	for formula in $(FORMULAS); do cd $$formula/src && make build && cd $(PWD_INITIAL); done
	./copy-bin-configs.sh "$(FORMULAS)"

test-local:
ifneq ("$(FORM)", "")
	@echo "Using form true: "  $(FORM_TO_UPPER)
	$(MAKE) bin FORMULAS=$(FORM)
	mkdir -p $(HOME)/.rit/formulas
	rm -rf $(HOME)/.rit/formulas/$(FORM)
	./unzip-bin-configs.sh
	cp -r formulas/* $(HOME)/.rit/formulas
	rm -rf formulas
else
	@echo "Use make test-local form=NAME_FORMULA for specific formula."
	@echo "form false: ALL FORMULAS"
	$(MAKE) bin
	rm -rf $(HOME)/.rit/formulas
	./unzip-bin-configs.sh
	mv formulas $(HOME)/.rit
endif
	mkdir -p $(HOME)/.rit/repo/local
	rm -rf $(HOME)/.rit/repo/local/tree.json
	cp tree/tree.json  $(HOME)/.rit/repo/local/tree.json

