#/usr/bin/env python
"""
xxxx
[para102]
        address 10.46.0.102
        use_node_name yes



"""


import sys

ip = sys.argv[1]

with open('/etc/munin/munin.conf', 'r+') as f:
  data=f.read()
  f.seek(0)
  data = re.sub(
            re.compile(r"(\[[A-Za-z0-9]+\])"
                       ""
                       "([\t ]+)address "+ip
                       "", re.MULTILINE), replaceWith, data
            )
  #print data
  f.write(data)
  f.truncate()
  f.close()