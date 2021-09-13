#!/usr/bin/python3
import os

from hello import hello

input1 = os.environ.get('SAMPLE_TEXT')
input2 = os.environ.get('SAMPLE_LIST')
input3 = os.environ.get('SAMPLE_BOOL')
hello.Run(input1, input2, input3)