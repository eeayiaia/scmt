#/usr/bin/env python
import sys
import re

ip = sys.argv[1]

with open('/etc/munin/munin.conf', 'r+') as f:
  data=f.read()
  f.seek(0)
  #find the node to be removed
  regex = re.compile(r"^(\[(.+)\])((\n\t(.+)(\s)(.+))*)(\n\t(address)(\s))"+ip+"((\n\t(.+)(\s)(.+))*)", re.MULTILINE)
  #replace node configuration with the empty string
  newconf = re.sub(regex, "", data)
  
  if newconf == data:
    print "no node could be found"
    sys.exit()
  
  f.write(newconf)
  f.truncate()
  f.close()
