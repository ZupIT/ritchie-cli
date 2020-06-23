#!/bin/sh

BETA_VERSION=$(expr $(curl -s https://commons-repo.ritchiecli.io/stable.txt| rev | cut -d . -f -1|rev) + 1)

if [ $? != 0 ]; then
        BETA_VERSION="5"
fi

echo "$VERSION_PLACEHOLDER" | sed "s/PLACEHOLDER/.pre.${BETA_VERSION}/"