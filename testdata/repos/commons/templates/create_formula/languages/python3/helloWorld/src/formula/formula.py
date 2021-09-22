#!/usr/bin/python3
from colored import fg, attr
from distutils.util import strtobool


def Run(input1, input2, input3, input4):
    print("Hello World!")
    print(f"{fg(2)}My name is {input1}.{attr(0)}")
    if strtobool(input2):
        print(f"{fg(3)}I've already created formulas using Ritchie.{attr(0)}")
    else:
        s = "I'm excited in creating new formulas using Ritchie."
        print(f"{fg(3)}'{s}'.{attr(0)}")
    print(f"{fg(1)}Today, I want to automate {input3}.{attr(0)}")
    print(f"{fg(3)}My secret is '{input4}'.{attr(0)}")
