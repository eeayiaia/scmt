#!/bin/bash


#Get script directory
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi

. "$DIR/../../../scripts.d/utils.sh" || exit 1
. "$DIR/../resources/config" || exit 1


#Create hadoop user & give sudo rights
useradd -m hadoop
passwd hadoop hadoop
sudo adduser hadoop sudo

#Create hadoop user group
sudo addgroup hadoop
sudo adduser --ingroup hadoop hadoop


#Check java version
if type -p java; then
	echo found java executable in PATH
	_java=java
elif [[ -n "$JAVA_HOME" ]] && [[ -x "$JAVA_HOME/bin/java" ]]; then
	echo "found java executable in $JAVA_HOME"
	_java="$JAVA_HOME/bin/java"
else
	echo "no java"
  echo "Installing java..."
	sudo apt-get update
	sudo apt-get install openjdk-8-jdk
	echo "Java installed."
fi

if [[ "$_java"]]; then
	version=$("$_java" - version 2>&1 | awk -F '"' '/version/ {print $3}')
	echo version "$version"
	if [[ "$version" == "1.8" ]]; then
		echo "version is 8"
	else
		echo "updating java to 1.8"
		sudo apt-get update
		sudo apt-get install openjdk-8-jdk
	fi
fi


#Dissable IPv6
#echo "net.ipv6.conf.all.disable_ipv6 = 1" >> /ect/systl.conf
#echo "net.ipv6.conf.default.diable_ipv6 = 1" >> /ect/systl.conf
#echo "net.ipv6.conf.lo.disable_ipv6 = 1" >> /ect/systl.conf


#Installing Hadoop below...
#Hadoop binarys needs to be shared in folder or downloaded before installation, else use line above
tar xfz hadoop-2.7.2.tar.gz
mv hadoop-2.7.2 /usr/local/hadoop


#Update .bashrc
echo "#Set Hadoop-releated enviroment variables" >> /home/hadoop/.bashrc
echo "export HADOOP_HOME=/home/hadoop" >> /home/hadoop/.bashrc
echo "#Set Java home " >> /home/hadoop/.bashrc
echo "export JAVA_HOME=/usr/lib/java-8-openjdk" >> /home/h/.bashrc


#Add nodes to hadoop system
#echo "" >> /ect/hosts
#echo "10.46.0.101			hadoop-master" >> /etc/hosts
#echo "10.46.0.102			hadoop-slave-1" >> /etc/hosts
#echo "10.46.0.103			hadoop-slave-2" >> /etc/hosts

#Configure Hadoop

#First we need to set the java home directory in hadoop-env.sh
cd  /usr/local/hadoop/etc/hadoop/
cp hadoop-env.sh hadoop-env-backup.sh
sed -e 's#export JAVA_HOME=${JAVA_HOME}#export JAVA_HOME=/usr/lib/jvm/java-1.8.0-openjdk-armf#g' 
sed -e 's#export HADOOP_OPTS="$HADOOP_OPTS -Djava.net.preferIPv4Stack=true"#export HADOOP_OPTS="$HADOOP_OPTS -Djava.net.preferIPv4Stack=true -Djava.library.path=$HADOOP_PREFIX/lib"#g'
echo "export HADOOP_IDENT_STRING=$USER
export HADOOP_COMMON_LIB_NATIVE_DIR=${HADOOP_PREFIX}/lib/native" >> /usr/local/hadoop/etc/hadoop/hadoop-env.sh

#Set Hadoop enviroment

#Config yarn-env.sh
cd /usr/local/hadoop/etc/hadoop/
cp yarn-env.sh yarn-env-backup.sh
echo "export HADOOP_CONF_LIB_NATIVE_DIR=${HADOOP_PREFIX:-"/lib/native"}
export HADOOP_OPTS="-Djava.LIBRARY.PATH=$HADOOP_PREFIX/lib"" >> /usr/local/hadoop/etc/hadoop/yarn-env.sh

#Config core-site.xml
cd /usr/local/hadoop/etc/hadoop/
cp core-site.xml core-site-backup.xml
sed -e 's#<configuration># #g' core-site.xml
sed -e 's#</configuration># #g' core-site.xml

echo "<configuration>
 <property>
   <name>hadoop.tmp.dir</name>
   <value>/app/hadoop/tmp</value>
   <description>A base for other temporary directories.</description>
 </property>

 <property>
   <name>fs.default.name</name>
   <value>hdfs://localhost:54310</value>
   <description>The name of the default file system.  A URI whose
   scheme and authority determine the FileSystem implementation.  The
   uri's scheme determines the config property (fs.SCHEME.impl) naming
   the FileSystem implementation class.  The uri's authority is used to
   determine the host, port, etc. for a filesystem.</description>
 </property>
</configuration>" >> /usr/local/hadoop/etc/hadoop/core-site.xml



#Config hdfs-site.xml
cd /usr/local/hadoop/etc/hadoop/
cp hdfs-site.xml hdfs-site-backup.xml
sed -e 's#<configuration># #g' hdfs-site.xml
sed -e 's#</configuration># #g' hdfs-site.xml

echo "<configuration>
<property>
  <name>dfs.replication</name>
  <value>3</value>
  <description>Default block replication.
  The actual number of replications can be specified when the file is created.
  The default is used if replication is not specified in create time.
  </description>
</property>
</configuration> " >> /usr/local/hadoop/etc/hadoop/hdfs-site.xml


#Config mapred-site.xml
cd /usr/local/hadoop/etc/hadoop/
cp mapred-site.xml mapred-site-backup.xml
sed -e 's#<configuration># #g' mapred-site.xml
sed -e 's#</configuration># #g' mapred-site.xml
echo "<configuration>
 <property>
  <name>mapreduce.framework.name</name>
  <value>yarn</value>
 </property>
</configuration>" >> /usr/local/hadoop/etc/hadoop/mapred-site.xml




