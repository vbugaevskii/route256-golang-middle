#!/usr/bin/env sh

set -e

SCRIPT_DIR=$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )
POSTGRES_URI="postgres://postgres:password@localhost:5434/loms?sslmode=disable&statement_cache_mode=describe"

goose -dir ${SCRIPT_DIR}/migrations postgres ${POSTGRES_URI} up
# goose -dir ${SCRIPT_DIR}/migrations postgres ${POSTGRES_URI} down
goose -dir ${SCRIPT_DIR}/migrations postgres ${POSTGRES_URI} status
