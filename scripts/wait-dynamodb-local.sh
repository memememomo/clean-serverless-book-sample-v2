#!/usr/bin/env bash

check() {
  echo "Wait for $1"
  try=0
  until curl -X POST --connect-timeout 10 --max-time 10 $1 &> /dev/null; do
    >&2 echo -n "."
    sleep 1
    try=$(expr $try + 1)
    if [ $try -ge 10 ]; then
      echo ""
      echo "Failed to wait for $1"
      break
    fi
  done
}

check $DYNAMO_LOCAL_ENDPOINT
