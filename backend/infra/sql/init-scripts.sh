#!/bin/bash
set -e

echo "Creating databases..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE DATABASE auth_db;
    CREATE DATABASE auction_db;
    CREATE DATABASE bidding_db;
    CREATE DATABASE notification_db;
EOSQL

echo "Initializing auth_db..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "auth_db" -f /docker-entrypoint-initdb.d/schemas/user_init.sql

echo "Initializing auction_db..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "auction_db" -f /docker-entrypoint-initdb.d/schemas/auction_init.sql

echo "Initializing bidding_db..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "bidding_db" -f /docker-entrypoint-initdb.d/schemas/bidding_init.sql

echo "Initializing notification_db..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "notification_db" -f /docker-entrypoint-initdb.d/schemas/notification_init.sql

echo "All databases initialized successfully."
