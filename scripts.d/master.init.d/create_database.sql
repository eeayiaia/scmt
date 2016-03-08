CREATE DATABASE cluster;
USE cluster;
CREATE TABLE devices (hwaddr CHAR(12) KEY, ip INT UNSIGNED NOT NULL, hname varchar(30) NOT NULL, username VARCHAR(20) NOT NULL, password VARCHAR(20)NOT NULL);
CREATE USER 'master'@'localhost' IDENTIFIED BY 'badpassword';
GRANT ALL PRIVILEGES ON cluster . * TO master@localhost;
FLUSH PRIVILEGES;
