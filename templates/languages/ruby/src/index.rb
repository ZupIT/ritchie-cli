#!/usr/bin/ruby

require 'bundler/setup'
require_relative 'formula/formula'

INPUT1 = ENV["SAMPLE_TEXT"]
INPUT2 = ENV["SAMPLE_LIST"]
INPUT3 = ENV["SAMPLE_BOOL"]

Run(INPUT1, INPUT2, INPUT3)
