#!/usr/bin/env bash

selfdir=$(dirname $0)
srcdir=$selfdir
pidfile=$selfdir/jkwxfucker_exec.pid
logdir=$selfdir/log
logfile=$logdir/exec_`date +%Y%m%d_%s_%N`.log


cd $selfdir
nohup ./dayexec > $logfile & echo $! > $pidfile
