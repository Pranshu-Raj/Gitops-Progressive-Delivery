#!/usr/bin/env bash
set -u
url=${1:?usage: load.sh <url> [interval-seconds]}
interval=${2:-0.2}

while true; do
  curl -s -o /dev/null -w '%{http_code}\n' "$url"
  sleep "$interval"
done
