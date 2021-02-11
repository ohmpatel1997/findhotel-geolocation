#!/bin/bash

echo -n "Enter Hostname: "
read -s rds_hostname

echo
echo -n "Enter Username: "
read -s rds_username

echo
echo -n "Enter Password: "
read -s rds_password

echo
echo -n "Enter Database Name: "
read -s rds_db_name

export RDS_HOSTNAME="${rds_hostname}"
export RDS_USERNAME="${rds_username}"
export RDS_PASSWORD="${rds_password}"
export RDS_DB_NAME="${rds_db_name}"

echo
