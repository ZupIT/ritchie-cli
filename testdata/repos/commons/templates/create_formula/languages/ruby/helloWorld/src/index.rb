#!/usr/bin/ruby

require 'bundler/setup'
require_relative 'formula/formula'

INPUT1 = ENV["RIT_INPUT_TEXT"]
INPUT2 = ENV["RIT_INPUT_BOOLEAN"]
INPUT3 = ENV["RIT_INPUT_LIST"]
INPUT4 = ENV["RIT_INPUT_PASSWORD"]

Run(INPUT1, INPUT2, INPUT3, INPUT4)
