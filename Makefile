
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

build:
	make gen
	go build -o ./bin/app ./cmd/app

run:
	make build
	./bin/app


wire-gen:
	wire ./internal/app/

gen:
	make wire-gen

migrate.up:
	migrate -path ./migrations -database 'postgres://$(PG_USER):$(PG_PASS)@$(PG_HOST):$(PG_PORT)/$(PG_NAME)?sslmode=disable' up

migrate.down:
	migrate -path ./migrations -database 'postgres://$(PG_USER):$(PG_PASS)@$(PG_HOST):$(PG_PORT)/$(PG_NAME)?sslmode=disable' down