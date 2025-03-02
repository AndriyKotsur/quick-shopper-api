#!/bin/bash

set -e
set -u

function create_database() {
	local database=$1
	echo "Creating database '$database' for user '$POSTGRES_USER'"
	psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
	    CREATE DATABASE $database;
	    GRANT ALL PRIVILEGES ON DATABASE $database TO $POSTGRES_USER;
EOSQL
}

if [ -n "$POSTGRES_DATABASE" ]; then
	echo "Database creation requested: $POSTGRES_DATABASE"
	create_database $POSTGRES_DATABASE
	echo "Database '$POSTGRES_DATABASE' created successfully"
else
	echo "No database creation requested"
fi
