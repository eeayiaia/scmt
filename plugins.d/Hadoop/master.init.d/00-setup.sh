#!/bin/bash

#Check java version
if type -p java; then
	echo found java executable in PATH
	_java=java
elif [[ -n "$JAVA_HOME" ]] && [[ -x "$JAVA_HOME/bin/java" ]]; then
	echo "found java executable in $JAVA_HOME"
	_java="$JAVA_HOME/bin/java"
else
	echo "no java"
fi

if [[ "$_java"]]; then
	version=$("$_java" - version 2>&1 | awk -F '"' '/version/ {print $2}')
	echo version "$version"
	if [[ "$version" > "1.5" ]]; then
		echo "version is greater than 1.5"
	else
		echo updating java to 1.8
		sudo add-apt-repository ppa:webupd8team/java
		sudo apt-get update
		sudo apt-get install oracle-java8-installer
	fi
fi


#Create hadoop user
sudo addgroup hadoop
sudo adduser --ingroup hadoop hduser

#Dissable IPv6
#echo "net.ipv6.conf.all.disable_ipv6 = 1" >> /ect/systl.conf
#echo "net.ipv6.conf.default.diable_ipv6 = 1" >> /ect/systl.conf
#echo "net.ipv6.conf.lo.disable_ipv6 = 1" >> /ect/systl.conf


#Installing Hadoop below...
sudo add-apt-repository ppa:hadoop-ubuntu/stable
sudo apt-get update && sudo apt-get upgrade
sudo apt-get install hadoop

#Update .bashrc
echo "#Set Hadoop-releated enviroment variables" >> /home/hduser/.bashrc
echo "export HADOOP_HOME=/home/hduser/hadoop" >> /home/hduser/.bashrc
echo "#Set Java home " >> /home/hduser/.bashrc
echo "export JAVA_HOME=/usr/lib/java-8-openjdk" >> /home/hduser/.bashrc


#Add nodes to hadoop system
echo "" >> /ect/hosts
echo "10.46.0.101			hadoop-master" >> /etc/hosts
echo "10.46.0.102			hadoop-slave-1" >> /etc/hosts
echo "10.46.0.103			hadoop-slave-2" >> /etc/hosts

#Configure Hadoop
cd /opt/hadoop/hadoop/
































