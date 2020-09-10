#!/usr/bin/python3
import os

from formula import formula

input1 = os.environ.get("SAMPLE_TEXT")
input2 = os.environ.get("SAMPLE_LIST")
input3 = os.environ.get("SAMPLE_BOOL")
formula.Run(input1, input2, input3)
