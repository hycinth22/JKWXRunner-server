#!/usr/bin/env bash

selfDir=$(dirname "$0")
srcDir=${selfDir}
pidFile=${selfDir}/jkwxfucker.pid
logDir=${selfDir}/data/logs/api
logfile=${logDir}/api_$(date +%Y%m%d_%s_%N).log
tmp_exec=$(mktemp)
lastCommit=$(git log -1 --format="%H")
lastCommitTime=$(git log -1 --format="%cd" --date=relative)" (now $(date "+%F %T"))"
ldflags="-X 'main.lastCommit=${lastCommit}' -X 'main.lastCommitTime=${lastCommitTime}'"

if [[ ! -d ${logDir} ]]; then
  mkdir -m 655 "${logDir}"
fi

echo "building."
go build -ldflags "${ldflags}" -o "${tmp_exec}" "${srcDir}"
if [[ $? -ne 0 ]]; then
  echo "build failed."
  exit $?
fi
echo "build ok. ""${tmp_exec}"

nohup "${tmp_exec}" >"${logfile}" 2>&1 &
echo $! >"${pidFile}"
