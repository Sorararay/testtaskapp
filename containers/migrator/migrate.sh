#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

if [ "$#" -ne 7 ]; then
  echo "usage: $0 <db_host> <db_port> <db_name> <username> <password> <command> <ssl_mode>"
  exit 1
fi

db_host=$1
db_port=$2
db_name=$3
db_username=$4

if [ -e "$5" ]; then
  db_password=`cat $5`
else
  db_password=$5
fi

command=$6
ssl_mode=$7


echo "Waiting for PostgreSQL to start..."
until PGPASSWORD=$db_password psql -U $db_username -h $db_host -p $db_port -lqt &> /dev/null; do
  >&2 echo "PostgreSQL is unavailable - sleeping"
  sleep 1
done
echo "PostgreSQL is up - executing command"

migrate -path ./history -database postgres://$db_username:$db_password@$db_host:$db_port/$db_name?sslmode=$ssl_mode $command
