#!/bin/sh

NIGHTLY_VERSION=$(expr $(curl -s https://commons-repo.ritchiecli.io/stable.txt| rev | cut -d . -f -1|rev) + 1)

echo "$VERSION_PLACEHOLDER" | sed "s/PLACEHOLDER/.nightly.${NIGHTLY_VERSION}/"