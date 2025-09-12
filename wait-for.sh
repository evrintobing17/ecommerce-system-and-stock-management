#!/usr/bin/env bash
# wait-for.sh host:port -- command args...

set -e

host="$1"
shift
cmd="$@"

until nc -z ${host%:*} ${host#*:}; do
  echo "Waiting for $host..."
  sleep 2
done

>&2 echo "$host is up - executing command"
exec $cmd
