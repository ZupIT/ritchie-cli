#!/bin/bash

if expr "$CIRCLE_BRANCH" : 'qa' >/dev/null; then
  export AWS_ACCESS_KEY_ID="$AWS_ACCESS_KEY_ID_QA"
  export AWS_SECRET_ACCESS_KEY="$AWS_SECRET_ACCESS_KEY_QA"

elif expr "$CIRCLE_BRANCH" : '^release-.*' >/dev/null; then
  export AWS_ACCESS_KEY_ID="$AWS_ACCESS_KEY_ID_PROD"
  export AWS_SECRET_ACCESS_KEY="$AWS_SECRET_ACCESS_KEY_PROD"

elif expr "$CIRCLE_BRANCH" : '^nightly' >/dev/null; then
  export AWS_ACCESS_KEY_ID="$AWS_ACCESS_KEY_ID_PROD"
  export AWS_SECRET_ACCESS_KEY="$AWS_SECRET_ACCESS_KEY_PROD"

elif expr "$CIRCLE_BRANCH" : '^beta' >/dev/null; then
  export AWS_ACCESS_KEY_ID="$AWS_ACCESS_KEY_ID_PROD"
  export AWS_SECRET_ACCESS_KEY="$AWS_SECRET_ACCESS_KEY_PROD"

else
  echo ""
fi
