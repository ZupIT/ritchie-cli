#!/bin/bash

function bucket {

    VERSION=$RELEASE_VERSION

    if expr "$VERSION" : '^[0-9]\+' >/dev/null; then
      echo "ritchie-7395046262137"
    else
      echo ""
    fi

}

function credentials {

    if expr "$CIRCLE_BRANCH" : '^release-.*' >/dev/null; then
      export AWS_ACCESS_KEY_ID="$AWS_ACCESS_KEY_ID_PROD"
      export AWS_SECRET_ACCESS_KEY="$AWS_SECRET_ACCESS_KEY_PROD"
      export METRIC_BASIC_USER="$METRIC_BASIC_USER_PROD"
      export METRIC_BASIC_PASS="$METRIC_BASIC_PASS_PROD"

    else
      echo ""
    fi

}

function gonna_release {

    DEPLOYED_VERSION=$(curl -s https://commons-repo.ritchiecli.io/stable.txt)
    DIFF_RESULT=$(git --no-pager log --oneline beta...${DEPLOYED_VERSION} 2>/dev/null)

    if [ -z "$DIFF_RESULT" ]
    then
          echo "ABORT"
    else
          echo "RELEASE"
    fi


}

function next_version {

    NEXT_VERSION=$(expr $(curl -s https://commons-repo.ritchiecli.io/stable.txt| rev | cut -d . -f -1|rev) + 1)
    echo "${VERSION_PLACEHOLDER//PLACEHOLDER/$NEXT_VERSION}"

}


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

    if expr "$CIRCLE_BRANCH" : '^release-.*' >/dev/null; then
      export RELEASE_VERSION=$(echo "$CIRCLE_BRANCH" | cut -d '-' -f 2-)
    else
      export RELEASE_VERSION=$(curl https://commons-repo.ritchiecli.io/stable.txt)
    fi

}

function caller {

   if expr "$1" : "bucket" >/dev/null; then
      version
      bucket
   elif expr "$1" : "credentials" >/dev/null; then
      credentials
   elif expr "$1" : "gonna_release" >/dev/null; then
      gonna_release
   elif expr "$1" : "next_version" >/dev/null; then
      next_version
   elif expr "$1" : "metric_server" >/dev/null; then
      version
      metric_server
   elif expr "$1" : "version" >/dev/null; then
      version
   else
     echo "Unable to process params"
     exit 1
   fi

}


caller "$1"
