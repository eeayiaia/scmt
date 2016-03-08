#!/bin/sh



MasterNodeName="Epa"

if [[ $EUID -ne 0 ]]; then
	echo "This installer must be run with root rights." 1>&2
	exit 100
fi

#Install dependencies
apt-get install -y apache2 apache2-utils libapache2-mod-fcgid libcgi-fast-perl

#Graph zooming packages
apt-get install -y libcgi-fast-perl libapache2-mod-fcgid

#Check that the zooming packages are installed and enabled

CGI_CHECK=`/usr/sbin/apachectl -M | grep -i fcgid_module`

#if no fcgid_module, try enabling it
if [[ -z $CGI_CHECK ]]; then
	a2enmod fcgid
fi
#TODO: what if it still does not work?

#Install Munin
apt-get install -y munin

#Add symbol link in conf-enabled for munin apache config.
ln -S /etc/munin/apache.conf /etc/apache2/conf-enabled/munin.conf

#Set master node name
sed -i 's/localhost\.localdomain/$MasterNodeName/g' /etc/munin/munin.conf

#If apache version >= 2.4 change munin config to be compatible.

apacheVersion=`apachectl -v | head -n 1 | cut -c24-26`
if [ `echo "$apacheVersion >= 2.4" | bc -l` ]; then
  echo "Current apache config with munin incompatible, fixing it."
  sed -i '/Options None/d' /etc/munin/apache.conf
  sed -i 's/Order allow,deny/Require all granted/g' /etc/munin/apache.conf
  sed -i 's/Allow from localhost 127\.0\.0\.0\/8 ::1/Options FollowSymLinks SymLinksIfOwnerMatch/g' /etc/munin/apache.conf
fi

service apache2 restart

#Get package including deps list url for debs
#apt-get --print-uris --yes install apache2 | grep ^\' | cut -d\' -f2
#not working if package already installed.
