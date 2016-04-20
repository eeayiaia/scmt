#!/bin/bash

# Input: MYSQL_PASSWORD

mysql -u root -p$MYSQL_PASSWORD < create_database.sql
