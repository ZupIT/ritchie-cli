#!/bin/sh

if expr "$VERSION" : '.*qa.*' >/dev/null; then
  echo "https://ritchie-metrics.itiaws.dev/metrics"
elif expr "$VERSION" : '^[0-9]\+' >/dev/null; then
  echo "https://ritchie-metrics.zup.io/metrics"
else
  echo ""
fi
