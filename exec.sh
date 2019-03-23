#!/usr/bin/env bash

params=$*
selfDir=$(dirname $0)
srcDir=${selfDir}
pidFile=${selfDir}/jkwxfucker_exec.pid
logDir=${selfDir}/log
logfile=${logDir}/exec_`date +%Y%m%d_%s_%N`.log


cd ${selfDir}
nohup ./dayexec ${params} > ${logfile} & echo $! > ${pidFile}
