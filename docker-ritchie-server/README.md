## Docker Compose Ritchie-Server
### Run
```
docker-compose up
```

### Server configuration
```
./data/server/file_config.json
```

### Cli using keycloak
```
rit init
✔ yes
Enter your organization:  keycloak|
✔ yes
URL of the server [http(s)://host]:  http://localhost:3000|
✔ yes
Enter your username:  user|
Enter your password:  ****|
Organization: ldap
Login successfully!
```
User: user

Password: admin

### Cli using Ldap
```
rit init
✔ yes
Enter your organization:  ldap|
✔ yes
URL of the server [http(s)://host]:  http://localhost:3000|
✔ yes
Enter your username:  user|
Enter your password:  ****|
Organization: ldap
Login successfully!
```
User: user

Password: user

