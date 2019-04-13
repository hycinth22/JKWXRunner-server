#!/usr/bin/env bash

params=$*
selfDir=$(dirname $0)
srcDir=${selfDir}/executor
pidFile=${selfDir}/jkwxfucker_exec.pid
logDir=${selfDir}/log
logfile=${logDir}/exec_`date +%Y%m%d_%s_%N`.log


cd ${selfDir}
nohup go run ${srcDir} ${params} > ${logfile} & echo $! > ${pidFile}
