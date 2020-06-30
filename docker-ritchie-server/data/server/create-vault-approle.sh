#!/bin/sh

C_DIR=/home/application
VAULT_HOST=$1

mkdir -p /tmp/vault

echo "Exporting vault vars..."
export VAULT_TOKEN="87e7784b-d598-44fe-8962-c7c345a11eed"
export VAULT_ADDR=$VAULT_HOST

apk --no-cache add jq curl
curl -fsSL https://releases.hashicorp.com/vault/1.3.0/vault_1.3.0_linux_amd64.zip -o /tmp/vault/vault_1.3.0_linux_amd64.zip

echo "Installing vault cli 1.3.0..."
rm -rf /tmp/vault/vault
unzip /tmp/vault/vault_1.3.0_linux_amd64.zip -d /tmp/vault/

/tmp/vault/vault secrets enable -path=ritchie/warmup generic
/tmp/vault/vault secrets enable -path=ritchie/credential generic
/tmp/vault/vault secrets enable -path=ritchie/transit transit

/tmp/vault/vault policy write ritchie_server_policy $C_DIR/ritchie_server_policy.hcl

/tmp/vault/vault auth enable approle
/tmp/vault/vault write auth/approle/role/ritchie_credential_role policies=ritchie_server_policy period=15s

# Create a key for transit (encrypt and decrypt)
/tmp/vault/vault write -f ritchie/transit/keys/ritchie_key

rm -f /tmp/vault/role-id.txt
rm -f /tmp/vault/secret-id.txt

role_response=$(/tmp/vault/vault read -format=json auth/approle/role/ritchie_credential_role/role-id)
echo "role_response $role_response"
role_id=$(echo $role_response | jq -j '.data.role_id')
echo "role_id: $role_id"
eval echo $role_id >> /tmp/vault/role-id.txt

secret_response=$(/tmp/vault/vault write -force -format=json auth/approle/role/ritchie_credential_role/secret-id)
echo "secret_response: $secret_response"
secret_id=$(echo $secret_response | jq -j '.data.secret_id')
echo "secret_id: $secret_id"
eval echo $secret_id >> /tmp/vault/secret-id.txt

unset VAULT_TOKEN