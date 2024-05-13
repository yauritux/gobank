#!/bin/bash
set -e

psql -v ON ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE EXTENSION "uuid-ossp";
    CREATE USER gobank;
    CREATE DATABASE gobank;
    GRANT ALL PRIVILEGES ON DATABASE gobank TO gobank;
EOSQL