#!/bin/sh
set -e

echo "Waiting for backend to become resolvable..."
# loop until backend resolves
until getent hosts backend >/dev/null 2>&1; do
  sleep 1
done

echo "Backend is resolvable. Starting Nginx..."
nginx -g "daemon off;"
