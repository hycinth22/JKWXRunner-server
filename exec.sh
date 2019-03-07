#!/bin/bash

selfdir=$(dirname $0)
srcdir=$selfdir
pidfile=$selfdir/jkwxfucker_exec.pid
logdir=$selfdir/log
logfile=$logdir/exec_`date +%Y%m%d_%s_%N`.log

# check running process
# we haven't implement it yet.
# 
# pid=$(cat $pidfile)
#if [ $pid -ne "" -a $pid -gt 0 ]; then
#	echo "wait pid $pid exit"
#	while [ -e /proc/ ]; do sleep 1; done
# fi'


cd $srcdir
nohup go run $srcdir/exec > $logfile & echo $! > $pidfile