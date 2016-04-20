#!/usr/bin/env python

import fileinput
import sys
import re
ip = sys.argv[1]
for line in fileinput.input("/etc/ganglia/gmetad.conf", inplace=True):
    if line[:11] == "data_source":
        e=True
        line = line.replace(" "+ip,"")
    print "%s" % (line),