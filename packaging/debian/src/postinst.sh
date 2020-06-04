#!/bin/sh -xe

mkdir ~/.rit

cat <<EOT > ~/.rit/server.json
{
  "organization": "zup",
  "url": "https://ritchie-server.zup.io"
}
EOT