#!/usr/bin/perl
use warnings;
use strict;
use Formula::Formula;

my $ritInputText = $ENV{'RIT_INPUT_TEXT'};
my $ritInputBoolean = $ENV{'RIT_INPUT_BOOLEAN'};
my $ritInputList = $ENV{'RIT_INPUT_LIST'};
my $ritInputPass = $ENV{'RIT_INPUT_PASSWORD'};

Formula::run($ritInputText, $ritInputBoolean, $ritInputList, $ritInputPass);
