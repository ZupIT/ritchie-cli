package tpl_php

const (
	Index = `<?php

	echo '123';

// $input1 = os.environ.get('SAMPLE_TEXT')
// $input2 = os.environ.get('SAMPLE_LIST')
// $input3 = os.environ.get('SAMPLE_BOOL')
// {{bin-name}}.Run(input1, input2, input3)

?>
`

	Makefile = `# Make Run PHP
BINARY_NAME={{bin-name}}.sh
BINARY_NAME_WINDOWS={{bin-name}}.bat
DIST=../dist
DIST_DIR=$(DIST)/commons/bin
build:
	mkdir -p $(DIST_DIR)
	cp index.php $(DIST_DIR) && cp -r {{bin-name}} Dockerfile set_umask.sh $(DIST_DIR)
	chmod +x $(DIST_DIR)/index.php
	echo 'php index.php' >> $(DIST_DIR)/$(BINARY_NAME_WINDOWS)`

	Dockerfile = `
FROM php:7.4-cli

COPY . .

RUN chmod +x set_umask.sh

WORKDIR /app

ENTRYPOINT ["/set_umask.sh"]

CMD ["php", "./index.php"]`

Run = `#!/bin/sh

php -f index.php

`

	File = `<?php

	function Run($input1, $input2, $input3) {
		echo "Hello World!";
		echo "You receive $input1 in text.";
		echo "You receive $input2 in list.";
		echo "You receive $input3 in boolean.";
	}
?>
`
)
