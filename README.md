# go-gin-xorm-starter

![CodeBuildBadge](https://codebuild.ap-northeast-1.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiYTFZRitkRDdWdzdscEZhaHN6R2VnVDBBaXpVNHprVFdPZGhDQVZOQVdxZ1J6SnBCeFl1ZmZncm9aZ2JYaTJkZG01VU00MUFzalhzSjU0ZGhDRlluNy84PSIsIml2UGFyYW1ldGVyU3BlYyI6ImJqS2h5cGpibDY4MzdpYi8iLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=master) [![Coverage Status](https://coveralls.io/repos/github/itomofumi/go-gin-xorm-starter/badge.svg)](https://coveralls.io/github/itomofumi/go-gin-xorm-starter)

Golang API starter using Gin and xorm.

## Requirements

- go (>= v1.10 recommended)
- dep
- docker & docker-compose
- nodemon for Live Reloading Development

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

### Install dependencies

```sh
go get
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

open <http://localhost:9080/db_structure.php?db=go-gin-xorm-starter>

Or, access Database using Adminer.

open <http://localhost:10080/?server=mysql&username=root&pass=password>

### Let's get fruits!!

Start API server.

```sh
make start
```

get fruits using `curl`.

```sh
curl http://localhost:3000/v1/fruits
```

### Add new fruit

post fruit using `curl`.

```sh
curl -X POST \
  -H 'Authorization:Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImlhdCI6MTUxNjIzOTAyMn0.hkDGuuaVbg2rBeEk3e97yUzl3Gp2UfD_hZO0dnjH6elS4WmxplQzXEXdOSvVaGFTxLpvwvTx11MT3PZzBUkoKR7WkGa76YdKiJGR-SZy7Zpdj6u1FdB9BGsIuvnfl0foX8En2JPV-EIA5Pm2fdy2hSGg1nzaPMekL8KeEJYjyi8' \
  -d '{"name":"Lemon","price":144}' \
  http://localhost:3000/v1/fruits
```

This sample JWT is generated [here](https://jwt.io/#debugger-io?token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImlhdCI6MTUxNjIzOTAyMn0.hkDGuuaVbg2rBeEk3e97yUzl3Gp2UfD_hZO0dnjH6elS4WmxplQzXEXdOSvVaGFTxLpvwvTx11MT3PZzBUkoKR7WkGa76YdKiJGR-SZy7Zpdj6u1FdB9BGsIuvnfl0foX8En2JPV-EIA5Pm2fdy2hSGg1nzaPMekL8KeEJYjyi8&publicKey=-----BEGIN%20PUBLIC%20KEY-----%0AMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDdlatRjRjogo3WojgGHFHYLugd%0AUWAY9iR3fy4arWNA1KoS8kVw33cJibXr8bvwUAUparCwlvdbH6dvEOfou0%2FgCFQs%0AHUfQrSDv%2BMuSUMAe8jzKE4qW%2BjK%2BxQU9a03GUnKHkkle%2BQ0pX%2Fg6jXZ7r1%2FxAK5D%0Ao2kQ%2BX5xK9cipRgEKwIDAQAB%0A-----END%20PUBLIC%20KEY-----).

### Shutdown

Stop Docker.

```sh
docker-compose down
```

## Build Docker image

Build.

```sh
docker build -t myimage/starter .
```

Run.

Use `-e` and `--env-file` option to give environment variables.

```sh
docker run -p 8888:80 --env-file .env -e DATABASE_HOST=mydbhost myimage/starter
```
