#!/usr/bin/python2
from colored import fg, attr
from distutils.util import strtobool


def Run(input1, input2, input3, input4):
    print "Hello World!"
    print "%sMy name is %s.%s" % (fg(2), input1, attr(0))

    if strtobool(input2):
        print "%sI've already created formulas using Ritchie.%s" % (fg(3), attr(0))
    else:
        print "%sI'm excited in creating new formulas using Ritchie.%s" % (
            fg(3),
            attr(0),
        )
    print "%sToday, I want to automate %s.%s" % (fg(1), input3, attr(0))
    print "%sMy secret is '%s'.%s" % (fg(3), input4, attr(0))
