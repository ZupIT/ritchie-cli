#!/bin/sh

NIGHTLY_VERSION=$(expr $(curl -s https://commons-repo.ritchiecli.io/stable.txt| sed 's/./&\n/g' | tac | sed -e :a -e 'N;s/\n//g;ta' | cut -d . -f -1| sed 's/./&\n/g' | tac | sed -e :a -e 'N;s/\n//g;ta') + 1)

if [ $? != 0 ]; then
        NIGHTLY_VERSION="4"
fi

echo "$VERSION_PLACEHOLDER" | sed "s/PLACEHOLDER/.nightly.${NIGHTLY_VERSION}/"