package template_php

const (
	Index = `<?php

include '{{bin-name}}/{{bin-name}}.php';

$input1 = getenv('SAMPLE_TEXT');
$input2 = getenv('SAMPLE_LIST');
$input3 = getenv('SAMPLE_BOOL');

Run($input1, $input2, $input3);

?>
`

	Makefile = `# Make Run PHP
BINARY_NAME_UNIX={{bin-name}}.sh
BINARY_NAME_WINDOWS={{bin-name}}.bat
DIST=../dist
DIST_DIR=$(DIST)/commons/bin
build:
	mkdir -p $(DIST_DIR)
	cp run_template $(BINARY_NAME_UNIX) && chmod +x $(BINARY_NAME_UNIX)
	echo 'php index.php' >> $(DIST_DIR)/$(BINARY_NAME_WINDOWS)

	cp -r . $(DIST_DIR)

	#Clean files
	rm $(BINARY_NAME_UNIX)`

	Dockerfile = `
FROM php:latest

COPY . .

RUN chmod +x set_umask.sh
RUN mkdir app
ENTRYPOINT ["/set_umask.sh"]
CMD ["php /index.php"]`

	Run = `#!/bin/sh

php -f index.php

`

	File = `<?php

	function Run($input1, $input2, $input3) {
		echo "Hello World! \n";
		echo "You receive $input1 in text. \n";
		echo "You receive $input2 in list. \n";
		echo "You receive $input3 in boolean. \n";
	}
?>
`
)
