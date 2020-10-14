#!/bin/sh

if expr "$VERSION" : '.*qa.*' >/dev/null; then
  echo "https://ritchie-metrics.devdennis.zup.io/v2/metrics"
elif expr "$VERSION" : '.*stg.*' >/dev/null; then
  echo "https://ritchie-metrics.stgdennis.zup.io/v2/metrics"
elif expr "$VERSION" : '^[0-9]\+' >/dev/null; then
  echo "https://ritchie-metrics.zup.io/v2/metrics"
else
  echo ""
fi
