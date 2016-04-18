CREATE DATABASE IF NOT EXISTS cluster;
USE cluster;
CREATE TABLE IF NOT EXISTS devices (hwaddr CHAR(12) KEY, id INT UNSIGNED AUTO_INCREMENT, port INT UNSIGNED, hname varchar(30), username VARCHAR(20), password VARCHAR(20), KEY `id` (`id`));
CREATE TABLE IF NOT EXISTS plugins (
    name VARCHAR(30) PRIMARY KEY,
    enabled BOOLEAN NOT NULL DEFAULT 0,
    installedOnMaster BOOLEAN NOT NULL DEFAULT 0
);
CREATE TABLE IF NOT EXISTS installedPlugins_slave (
    hwaddr CHAR(12) NOT NULL,
    plugin VARCHAR(30) NOT NULL,
    FOREIGN KEY (hwaddr) REFERENCES devices(hwaddr),
    FOREIGN KEY (plugin) REFERENCES plugins(name),
    UNIQUE (hwaddr, plugin)
);

GRANT ALL PRIVILEGES ON cluster . * TO master@localhost IDENTIFIED BY 'badpassword';
FLUSH PRIVILEGES;
