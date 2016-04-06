CREATE DATABASE IF NOT EXISTS cluster;
USE cluster;
CREATE TABLE IF NOT EXISTS devices (hwaddr CHAR(12) KEY, ip INT UNSIGNED NOT NULL, port INT UNSIGNED, hname varchar(30) NOT NULL, username VARCHAR(20) NOT NULL, password VARCHAR(20)NOT NULL);
CREATE TABLE IF NOT EXISTS plugins (
    name VARCHAR(30) PRIMARY KEY,
    enabled BOOLEAN NOT NULL DEFAULT 0
);
CREATE TABLE IF NOT EXISTS pluginsInstalledOn (
    hwaddr CHAR(12) NOT NULL,
    plugin VARCHAR(30) NOT NULL,
    FOREIGN KEY (hwaddr) REFERENCES devices(hwaddr),
    FOREIGN KEY (plugin) REFERENCES plugins(name),
    UNIQUE (hwaddr, plugin)
);

GRANT ALL PRIVILEGES ON cluster . * TO master@localhost IDENTIFIED BY 'badpassword';
FLUSH PRIVILEGES;
