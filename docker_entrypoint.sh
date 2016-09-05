#!/bin/sh
set -e

if [ "$1" = "pumba" ]; then
  if [ -S /var/run/docker.sock ]; then
    chown -R core:core /var/run/docker.sock
  fi
  exec gosu core:core "$@"
fi

exec "$@"
