#!/bin/sh

SCRIPT_DIR=`dirname $0`
ENV_FILE=$SCRIPT_DIR/../.env go run $SCRIPT_DIR/init.go $SCRIPT_DIR/db.sql
