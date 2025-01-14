#!/bin/bash
go mod tidy

docker-compose down

docker compose up -d --build

docker system prune -a -f