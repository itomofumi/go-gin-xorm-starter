# go-gin-xorm-starter

Golang API starter using Gin and xorm.

## Develop with local Database

### Start

Run Docker.

```sh
docker-compose up -d
```

Init Database.

```sh
sh ./fixtures/init_db.sh
```

Access Database using phpMyAdmin.

open <http://localhost:8080/db_structure.php?db=go_gin_xorm_starter>

Or, access Database using Adminer.

open <http://localhost:10080/?server=mysql&username=root&pass=password>

### Shutdown

Stop Docker.

```sh
docker-compose down
```
