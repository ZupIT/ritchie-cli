#!/bin/sh

/home/application/wait-for-it.sh "vault:8200" && echo "vault is up"
/home/application/wait-for-it.sh "keycloak:8080" && echo "keycloak is up"
/home/application/wait-for-it.sh "ldap:389" && echo "ldap is up"

/home/application/create-vault-approle.sh http://vault:8200

export VAULT_ADDR=http://vault:8200
export VAULT_AUTHENTICATION=APPROLE
export VAULT_ROLE_ID=$(cat /tmp/vault/role-id.txt)
export VAULT_SECRET_ID=$(cat /tmp/vault/secret-id.txt)
export FILE_CONFIG="/home/application/file_config.json"

./server/ritchie-server