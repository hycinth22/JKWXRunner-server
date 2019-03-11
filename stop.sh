#!/usr/bin/env bash

selfdir=$(dirname $0)
pidfile=$selfdir/jkwxfucker.pid

if [ -f $pidfile ]; then
	kill -9 `cat $pidfile`
fi
