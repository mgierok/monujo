#!/bin/bash

echo -n "Database hostname: "
read DB_HOSTNAME
echo -n "Database username: "
read DB_USERNAME
echo -n "Database password: "
read DB_PASSWORD
echo -n "Database name: "
read DB_NAME
echo -n "Path to quotes: "
read QUOTES_PATH

echo -n "" > conf/db.conf
echo "DB_HOSTNAME=$DB_HOSTNAME" >> conf/db.conf
echo "DB_USERNAME=$DB_USERNAME" >> conf/db.conf
echo "DB_PASSWORD=$DB_PASSWORD" >> conf/db.conf
echo "DB_NAME=$DB_NAME" >> conf/db.conf
echo "QUOTES_PATH=$QUOTES_PATH" >> conf/db.conf

