#!/usr/bin/env python

import re,sys

print "Python regex replace called."

#Replaces section in config file of style:
#sectionName { 
#    whatever
#     ....
#    whatever
#}
def replace(sectionName, textToMatch, replaceWith):
    return re.sub(
            re.compile("(?<="+sectionName+" \{ )[^\}]*", re.MULTILINE), replaceWith, textToMatch
            )

if len(sys.argv) != 4:
  sys.stderr.write("Param 1: content for cluster{ } section gmond.conf.\n" \
                   "Param 2: content for udp_send_channel{ } section gmond.conf.\n" \
                   "Param 3: content for data_source line in gmetad.conf.\n")
  exit(1)

with open('/etc/ganglia/gmond.conf', 'r+') as f:
  print "Editing gmond.conf..."
  data=f.read()
  f.seek(0)
  data = replace("cluster", data, sys.argv[1])
  data = replace("udp\_send\_channel", data, sys.argv[2])
  #print data
  f.write(data)
  f.truncate()
  f.close()

with open('/etc/ganglia/gmetad.conf', 'r+') as f:
  print "Editing gmetad.conf..."
  data=f.read()
  f.seek(0)
  #data = re.sub("data\_source \"my cluster\" localhost\n", sys.argv[3], data)
  #print data
  f.write(data)
  f.truncate()
  f.close()