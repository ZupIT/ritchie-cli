#!/usr/bin/ruby

require 'colorize'

def Run(input1, input2, input3, input4)
    puts "Hello World!"
    puts "My name is #{input1}.".green
    if input2 == "true" then
        puts "I've already created formulas using Ritchie.".blue
     else
        puts "I`m excited in creating new formulas using Ritchie.".red
     end
    puts "Today, I want to automate #{input3}.".yellow
    puts "My secret is #{input4}.".cyan
end
