#!/usr/bin/env bash

#nohup mysqld -u root &
#echo | ls /tmp/data/
unzip -j /tmp/data/data.zip '*.json' -d data
#(echo y | nohup mysqld -u root) &
nohup mysqld -u root &

sleep 10
#mysqladmin -u root status

#mysql -u root < echo "SET PASSWORD FOR 'root'@'localhost' = PASSWORD(1234);"
mysql -u root < go/src/highloadcup2017/storage/scheme.sql

./go/bin/highloadcup2017 80 data