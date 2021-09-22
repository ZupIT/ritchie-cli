package Formula;
use warnings;
use strict;
use Term::ANSIColor;

sub run {
	my ($ritInputText, $ritInputBoolean, $ritInputList, $ritInputPass) = @_;

	print "Hello World\n";
	print colored("My name is $ritInputText\n", "green");
	if($ritInputBoolean eq "true") {
		print colored("I've already created formulas using Ritchie.\n", "blue");
	} else {
		print colored("I'm excited in creating new formulas using Ritchie.\n", "red");
	}
	print colored("Today, I want to automate $ritInputList.\n", "yellow");
	print colored("My secret is $ritInputPass.\n", "cyan");
}

1;
