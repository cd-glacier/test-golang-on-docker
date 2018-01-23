#!/bin/bash -eu

mysql=( mysql --protocol=socket -uroot -p"${MYSQL_ROOT_PASSWORD}" )

"${mysql[@]}" <<-EOSQL
    CREATE DATABASE IF NOT EXISTS kyuko;
EOSQL

mysql -u root -ppassword kyuko < /docker-entrypoint-initdb.d/dump.sql_
