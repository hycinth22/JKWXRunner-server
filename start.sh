#!/usr/bin/env bash

selfdir=$(dirname $0)
srcdir=$selfdir
pidfile=$selfdir/jkwxfucker.pid
logdir=$selfdir/log
logfile=$logdir/api_`date +%Y%m%d_%s_%N`.log
tmp_exec=$(mktemp)

if [ ! -d $logdir ]; then
	mkdir -m 655 $logdir
fi

echo "building."
go build -o $tmp_exec $srcdir
if [ $? -ne 0 ]; then
	echo "build failed."
	exit $?
fi
echo "build ok. "$tmp_exec

nohup $tmp_exec > $logfile & echo $! > $pidfile
