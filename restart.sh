#!/usr/bin/env bash

selfdir=$(dirname $0)

$selfdir/stop.sh && $selfdir/start.sh
