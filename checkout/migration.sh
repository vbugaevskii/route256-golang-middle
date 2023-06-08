#!/usr/bin/env sh

set -e

SCRIPT_DIR=$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )

goose -dir ${SCRIPT_DIR}/migrations postgres "postgres://postgres:password@localhost:5433/checkout?sslmode=disable" up
# goose -dir ${SCRIPT_DIR}/migrations postgres "postgres://postgres:password@localhost:5433/checkout?sslmode=disable" down
goose -dir ${SCRIPT_DIR}/migrations postgres "postgres://postgres:password@localhost:5433/checkout?sslmode=disable" status
