#!/usr/bin/env bash

selfDir=$(dirname $0)

${selfDir}/stop.sh && ${selfDir}/start.sh
