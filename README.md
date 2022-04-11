# Desc
Tool does the periodic check of the crypto pairs api

## Features
- Makefile
- Async job (redis used)
- Support docker and docker-compose
- Redis used for the job sync
- Almost all is configurable

## Requiremetns
- Docker
- Docker compose
- Go 1.18 (generic used) in case of local build

## Installation
```sh
make build
make make run  
```

## TODO:
- Support WS
- Fix comments
