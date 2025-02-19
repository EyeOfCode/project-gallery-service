# Pre test gallery

## Can read api on swagger http://{host}:{port}/swagger

## stack

- docker
- docker-compose
- go fiber
- mongodb
- swagger
- air
- rate limit
- golangci-lint

## setup

- install go and setup path
- install docker and docker-compose

### create file .air.toml

```
root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o ./tmp/main.exe ./cmd/api"
bin = "tmp/main.exe"
full_bin = "./tmp/main.exe"
include_ext = ["go"]
exclude_dir = ["tmp", "mongodb_data"]
delay = 1000

[screen]
clear_on_rebuild = true
```

- init project $go mod init example-go-project
- init package $go mod tidy
- cp .env.example .env
- init swagger $swag init -g cmd/api/main.go
- build $go build cmd/api/main.go
- format $gofmt -w .
- lint $golangci-lint run

## how to use

- run $docker-compose up -d --build (init project or db)
- run app $go run cmd/api/main.go or use $air (air is build and compiler follow code change)

## Feature

- upload image
