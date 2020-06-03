package tpl_node

const (
	Index = `const run = require("./{{bin-name}}/{{bin-name}}")

const INPUT1 = process.env.SAMPLE_TEXT
const INPUT2 = process.env.SAMPLE_LIST
const INPUT3 = process.env.SAMPLE_BOOL

run(INPUT1, INPUT2, INPUT3)`

	Makefile = `# Make Run Node
BINARY_NAME_UNIX={{bin-name}}.sh
BINARY_NAME_WINDOWS={{bin-name}}.bat
DIST=../dist
DIST_DIR=$(DIST)/commons/bin

build:
	mkdir -p $(DIST_DIR)
	cp run_template $(BINARY_NAME_UNIX) && chmod +x $(BINARY_NAME_UNIX)
	echo 'node index.js' >> $(DIST_DIR)/$(BINARY_NAME_WINDOWS)

	cp -r $(BINARY_NAME_UNIX) index.js package.json {{bin-name}} $(DIST_DIR) && cp Dockerfile $(DIST_DIR)

	#Clean files
	rm $(BINARY_NAME_UNIX)`
	Dockerfile = `
FROM node:10

WORKDIR /app

COPY . .

RUN chmod +x set_umask.sh
ENTRYPOINT ["/app/set_umask.sh"]
CMD ["node index.js"]`

	Run = `#!/bin/sh
node index.js`

	PackageJson = `{
  "name": "src",
  "version": "1.0.0",
  "description": "Sample formula in node",
  "main": "index.js",
  "scripts": {
    "test": "echo \"Error: no test specified\" && exit 1"
  },
  "author": "Dennis.Ritchie",
  "license": "ISC"
}`

	File = `function Run(input1, input2, input3) {
    console.log("Hello World!")
    console.log("You receive "+ input1 +" in text.");
    console.log("You receive "+ input2 +" in list.");
    console.log("You receive "+ input3 +" in boolean.");
}

const {{bin-name}} = Run
module.exports = {{bin-name}}`
)
