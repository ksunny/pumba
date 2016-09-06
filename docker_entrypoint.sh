#!/bin/sh
set -e

if [ "$1" = "pumba" ]; then
#  if [ -S /var/run/docker.sock ]; then
#    chown -R root:root /var/run/docker.sock
#  fi
  exec gosu root:root "$@"
fi

exec "$@"
