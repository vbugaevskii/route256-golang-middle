#!/usr/bin/env sh

set -e

SCRIPT_DIR=$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )

goose -dir ${SCRIPT_DIR}/migrations postgres "postgres://postgres:password@localhost:5434/loms?sslmode=disable" up
# goose -dir ${SCRIPT_DIR}/migrations postgres "postgres://postgres:password@localhost:5434/loms?sslmode=disable" down
goose -dir ${SCRIPT_DIR}/migrations postgres "postgres://postgres:password@localhost:5434/loms?sslmode=disable" status
