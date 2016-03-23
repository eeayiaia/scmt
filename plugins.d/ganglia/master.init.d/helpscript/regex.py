#!/usr/bin/env python

import re,sys

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

if len(sys.argv)==4 and sys.argv[1]=="gmond":
  with open('/etc/ganglia/gmond.conf', 'r+') as f:
    print "Editing gmond.conf..."
    data=f.read()
    f.seek(0)
    data = replace(sys.argv[2], data, sys.argv[3])
    f.write(data)
    f.truncate()
    f.close()

elif len(sys.argv)==3 and sys.argv[1]=="gmetad":
  with open('/etc/ganglia/gmetad.conf', 'r+') as f:
    print "Editing gmetad.conf..."
    data=f.read()
    f.seek(0)
    data = re.sub("data\_source \"my cluster\" localhost\n", sys.argv[2], data)
    f.write(data)
    f.truncate()
    f.close()

else:
  print sys.argv[0] + " gmond section content"
  print sys.argv[0] + " gmetad data_source...."
  exit(0)
