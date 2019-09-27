#!/usr/bin/env bash

selfDir=$(dirname "$0")
pidFile=${selfDir}/jkwxfucker.pid

if [[ -f ${pidFile} ]]; then
  kill -9 "$(cat "${pidFile}")"
fi
