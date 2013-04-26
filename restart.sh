#!/bin/bash
BASE_DIR=`dirname $0`
BASE_DIR=`readlink -f $BASE_DIR`
pkill server
cd $BASE_DIR
nohup ./bin/server 2>$BASE_DIR/logs/error.log 1>$BASE_DIR/logs/server.log &
