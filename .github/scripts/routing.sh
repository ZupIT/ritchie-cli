#!/bin/bash

function metric_server {

    VERSION=$RELEASE_VERSION

    if expr "$VERSION" : '.*qa.*' >/dev/null; then
      echo "https://ritchie-metrics.devdennis.zup.io/v2/metrics"
    elif expr "$VERSION" : '.*stg.*' >/dev/null; then
      echo "https://ritchie-metrics.stgdennis.zup.io/v2/metrics"
    elif expr "$VERSION" : '^[0-9]\+' >/dev/null; then
      echo "https://ritchie-metrics.prddennis.zup.io/v2/metrics"
    else
      echo ""
    fi

}

function version {
      export RELEASE_VERSION=$(curl https://commons-repo.ritchiecli.io/stable.txt)
}

function caller {

   if expr "$1" : "metric_server" >/dev/null; then
      version
      metric_server
   else
     echo "Unable to process params"
     exit 1
   fi

}

caller "$1"
