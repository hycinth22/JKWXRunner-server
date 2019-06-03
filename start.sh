#!/usr/bin/env bash

selfDir=$(dirname $0)
srcDir=${selfDir}
pidFile=${selfDir}/jkwxfucker.pid
logDir=${selfDir}/data/logs/api
logfile=${logDir}/api_`date +%Y%m%d_%s_%N`.log
tmp_exec=$(mktemp)

if [[ ! -d ${logDir} ]]; then
	mkdir -m 655 ${logDir}
fi

echo "building."
go build -o ${tmp_exec} ${srcDir}
if [[ $? -ne 0 ]]; then
	echo "build failed."
	exit $?
fi
echo "build ok. "${tmp_exec}

nohup ${tmp_exec} > ${logfile} 2>&1 & echo $! > ${pidFile}
