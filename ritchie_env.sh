#!/bin/sh

if   expr "$VERSION" : 'qa-*' >/dev/null; then echo "qa"
elif expr "$VERSION" : '^[0-9]\+' >/dev/null; then echo "prod"
else echo ""
fi
