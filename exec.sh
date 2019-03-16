#!/usr/bin/env bash

params=$*
selfdir=$(dirname $0)
srcdir=$selfdir
pidfile=$selfdir/jkwxfucker_exec.pid
logdir=$selfdir/log
logfile=$logdir/exec_`date +%Y%m%d_%s_%N`.log


cd $selfdir
nohup ./dayexec $params > $logfile & echo $! > $pidfile
