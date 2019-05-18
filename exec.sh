#!/usr/bin/env bash

params=$*
selfDirRel=$(dirname "$0")
selfDir=$(readlink -f "$selfDirRel")
srcDir=${selfDir}/executor
pidFile=${selfDir}/jkwxfucker_exec.pid
logDir=${selfDir}/data/logs/exec
logfile=${logDir}/exec_`date +%Y%m%d_%s_%N`.log

echo Dir: ${selfDir}
cd ${selfDir}
nohup go run ${srcDir} ${params} > ${logfile} & echo $! > ${pidFile}
