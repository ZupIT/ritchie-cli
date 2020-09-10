#!/usr/bin/python3
from colored import fg, attr


def Run(input1, input2, input3):
    print("Hello World!")
    print("%sYou receive {} in text.%s".format(input1) % (fg(2), attr(0)))
    print("%sYou receive {} in list.%s".format(input2) % (fg(1), attr(0)))
    print("%sYou receive {} in boolean.%s".format(input3) % (fg(3), attr(0)))
