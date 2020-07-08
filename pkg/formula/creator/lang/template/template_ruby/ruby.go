package template_ruby

const (
	Index = `#!/usr/bin/ruby
require "./{{bin-name}}/{{bin-name}}"

INPUT1 = ENV["SAMPLE_TEXT"]
INPUT2 = ENV["SAMPLE_LIST"]
INPUT3 = ENV["SAMPLE_BOOL"]

Run(INPUT1, INPUT2, INPUT3)`

	Makefile = `# Make Run Ruby
BINARY_NAME_UNIX={{bin-name}}.sh
BINARY_NAME_WINDOWS={{bin-name}}.bat
DIST=../dist
DIST_DIR=$(DIST)/commons/bin

build:
	mkdir -p $(DIST_DIR)
	cp run_template $(BINARY_NAME_UNIX) && chmod +x $(BINARY_NAME_UNIX)
	echo './index.rb' >> $(DIST_DIR)/$(BINARY_NAME_WINDOWS)

	cp -r . $(DIST_DIR)

	#Clean files
	rm $(BINARY_NAME_UNIX)`

	Dockerfile = `FROM ruby:2.6

COPY . .

RUN chmod +x set_umask.sh

WORKDIR /app

ENTRYPOINT ["/set_umask.sh"]
CMD ["./index.rb"]
`

	Run = `#!/bin/sh
bundle install
./index.rb`

	Gemfile = `source 'https://rubygems.org' do
	# gem 'nokogiri'
	# Other gems here
end`

	File = `def Run(input1, input2, input3)
    puts "Hello World!"
    puts "You receive "+ input1 +" in text."
    puts "You receive "+ input2 +" in list."
    puts "You receive "+ input3 +" in boolean."
end`
)
