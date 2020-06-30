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

	Makefile = `# Make Run PHP
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

	WindowsBuild = `:: Php parameters
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
