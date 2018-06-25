# go-gin-xorm-starter

Golang API starter using Gin and xorm.

## Requirements

* go (>= v1.10 recommended)
* docker & docker-compose
* nodemon for Live Reloading Development

## Installation

### Install `nodemon` command

```sh
npm i -g nodemon

# or if you prefer yarn

yarn global add nodemon
```

### Setup `.env` File

Copy `.env.example` to `.env`.

```sh
cp .env.example .env
```

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

### Get Fruits

```sh
curl http://localhost:3000/v1/fruits
```

### Shutdown

Stop Docker.

```sh
docker-compose down
```
