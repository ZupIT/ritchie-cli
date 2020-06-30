#!/bin/sh

TIMEOUT=120

wait_for() {
  for i in `seq $TIMEOUT` ; do
    nc -z "$HOST" "$PORT" > /dev/null 2>&1
    
    result=$?
    if [ $result -eq 0 ] ; then
      echo "$HOST is running"
      exit 0
    fi
    sleep 1
  done
  echo "Operation timed out: $HOST" >&2
  exit 1
}

HOST=$(printf "%s\n" "$1"| cut -d : -f 1)
PORT=$(printf "%s\n" "$1"| cut -d : -f 2)

wait_for "$@"