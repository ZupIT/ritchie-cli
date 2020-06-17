#!/bin/sh

DEPLOYED_VERSION=$(curl -s https://commons-repo.ritchiecli.io/stable.txt)

VERSION_TO_CHECK_AGAINST=$($VERSION_PLACEHOLDER | sed "s/PLACEHOLDER//")

if [ "$DEPLOYED_VERSION" == "$VERSION_TO_CHECK_AGAINST" ]; then
    echo "RELEASE"
fi

echo "ABORT"