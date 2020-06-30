path "auth/token/*" {
  capabilities = [ "create", "read", "update", "delete", "list", "sudo" ]
}

path "auth/approle/login" {
  capabilities = [ "create", "read" ]
}

path "auth/approle/role/ritchie_credential_role/role-id" {
  capabilities = [ "read" ]
}

path "auth/approle/role/ritchie_credential_role/secret-id" {
  capabilities = ["create", "read", "update"]
}

path "secret/*" {
  capabilities = ["create", "read"]
}

path "ritchie/warmup/*" {
  capabilities = ["read","create","update"]
}

path "ritchie/credential/*" {
  capabilities = ["create", "update", "delete", "read", "list"]
}

path "ritchie/transit/encrypt/*" {
  capabilities = ["update"]
}

path "ritchie/transit/decrypt/*" {
  capabilities = ["update"]
}
