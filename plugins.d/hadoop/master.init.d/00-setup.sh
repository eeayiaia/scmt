#!/bin/bash


#Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../../scripts.d/utils.sh" || exit 1
. "$DIR/../resources/config" || exit 1


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
wget http://www.motorlogy.com/apache/hadoop/common/current/hadoop-2.3.0.tar.gz
tar xfz hadoop-2.3.0.tar.gz
mv hadoop-2.3.0 /usr/local/hadoop


#Update .bashrc
echo "#Set Hadoop-releated enviroment variables" >> /home/hduser/.bashrc
echo "export HADOOP_HOME=/home/hduser/hadoop" >> /home/hduser/.bashrc
echo "#Set Java home " >> /home/hduser/.bashrc
echo "export JAVA_HOME=/usr/lib/java-8-openjdk" >> /home/hduser/.bashrc


#Add nodes to hadoop system
#echo "" >> /ect/hosts
#echo "10.46.0.101			hadoop-master" >> /etc/hosts
#echo "10.46.0.102			hadoop-slave-1" >> /etc/hosts
#echo "10.46.0.103			hadoop-slave-2" >> /etc/hosts

#Configure Hadoop

#First we need to set the java home directory in hadoop-env.sh
sed 's/export JAVA_HOME=${JAVA_HOME}/export JAVA_HOME=${/usr/lib/jvm/java-1.7.0-openjdk-armf}' 
sed 's/export HADOOP_OPTS="$HADOOP_OPTS -Djava.net.preferIPv4Stack=true"/export HADOOP_OPTS="$HADOOP_OPTS -Djava.net.preferIPv4Stack=true -Djava.library.path=$HADOOP_PREFIX/lib"'
echo "export HADOOP_IDENT_STRING=$USER
export HADOOP_COMMON_LIB_NATIVE_DIR=${HADOOP_PREFIX}/lib/native" >> /usr/local/hadoop/etc/hadoop/hadoop-env.sh

#Set Hadoop enviroment

#Config yarn-env.sh
echo "export HADOOP_CONF_LIB_NATIVE_DIR=${HADOOP_PREFIX:-"/lib/native"}
export HADOOP_OPTS="-Djava.LIBRARY.PATH=$HADOOP_PREFIX/lib"" >> /usr/local/hadoop/etc/hadoop/yarn-env.sh

#Config core-site.xml




#Config hdfs-site.xml



#Config mapred-site.xml



#Config hadoop-env.sh
#Set this:  export JAVA_HOME=/usr/lib/jvm/java--openjdk-armf/jre/bin/jav export HADOOP_OPTS=-Djava.net.preferIPv4Stack=true export HADOOP_CONF_DIR=/opt/hadoop/hadoop/conf


































/usr/local/hadoop/etc/hadoop/hadoop-env.sh
