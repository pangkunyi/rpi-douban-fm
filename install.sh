#!/bin/bash
BASE_DIR=`dirname $0`
cd $BASE_DIR
export GOPATH=`pwd`
go install server
