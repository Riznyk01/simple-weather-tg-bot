#!/bin/sh
# wait-for-postgres.sh

set -e

host="$1"
shift
cmd="$@"

if [ -n "$DB_PASSWORD_FILE" ]; then
  # If POSTGRES_PASSWORD_FILE is set, read the content of the file and assign it to DB_PASSWORD
  DB_PASSWORD=$(cat "$DB_PASSWORD_FILE")
fi
until PGPASSWORD="$DB_PASSWORD" psql -h "$host" -U "postgres" -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - executing command"
exec $cmd