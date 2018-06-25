#!/bin/sh

SCRIPT_DIR=`dirname $0`
ENV_FILE=$SCRIPT_DIR/../.env go run fixtures/init.go fixtures/db.sql
