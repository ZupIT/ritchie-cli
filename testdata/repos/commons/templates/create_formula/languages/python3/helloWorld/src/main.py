#!/usr/bin/python3
import os

from formula import formula


input1 = os.environ.get("RIT_INPUT_TEXT")
input2 = os.environ.get("RIT_INPUT_BOOLEAN")
input3 = os.environ.get("RIT_INPUT_LIST")
input4 = os.environ.get("RIT_INPUT_PASSWORD")
formula.Run(input1, input2, input3, input4)
