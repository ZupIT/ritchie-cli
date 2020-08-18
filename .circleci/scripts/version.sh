#!/bin/sh

if expr "$CIRCLE_BRANCH" : 'qa' >/dev/null; then
  export RELEASE_VERSION="2.0.0-qa"
elif expr "$CIRCLE_BRANCH" : '.*beta.*' >/dev/null; then
  BETA_VERSION=$(expr $(curl -s https://commons-repo.ritchiecli.io/beta.txt| rev | cut -d . -f -1|rev) + 1)
  export RELEASE_VERSION=$(echo "$VERSION_PLACEHOLDER" | sed "s/PLACEHOLDER/.pre.${BETA_VERSION}/")
elif expr "$CIRCLE_BRANCH" : '^release-.*' >/dev/null; then
  export RELEASE_VERSION=$(echo "$CIRCLE_BRANCH" | cut -d '-' -f 2-)
elif expr "$CIRCLE_BRANCH" : '^nightly' >/dev/null; then
  NIGHTLY_VERSION=$(expr $(curl -s https://commons-repo.ritchiecli.io/stable.txt| rev | cut -d . -f -1|rev) + 1)
  export RELEASE_VERSION="$(echo "$VERSION_PLACEHOLDER" | sed "s/PLACEHOLDER/.nightly.${NIGHTLY_VERSION}/")"
else
  export RELEASE_VERSION=$(curl https://commons-repo.ritchiecli.io/stable.txt)
fi