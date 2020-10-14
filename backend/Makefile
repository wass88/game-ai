GO = go build
GO_ALPINE = GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

default: build

target/server: cmd/server/main.go  $(wildcard lib/server/*) 
	$(GO) -o $@ cmd/server/main.go

target/playout: cmd/playout/main.go $(wildcard lib/playout/*) $(wildcard lib/game/*) $(wildcard lib/game/*/*)
	$(GO) -o $@ cmd/playout/main.go

target/container: cmd/container/main.go $(wildcard lib/container/*) $(wildcard lib/docker/*)
	$(GO) -o $@ cmd/container/main.go

target/kick: cmd/kick/main.go $(wildcard lib/server/*)
	$(GO) -o $@ cmd/kick/main.go

target/server.alpine: cmd/server/main.go $(wildcard lib/server/*)
	$(GO_ALPINE) -o $@ cmd/server/main.go

docker-dev: target/server.alpine
	docker-compose build --no-cache
	docker-compose up 

db:
	docker-compose up -d db

dev: target/server
	CONF_FILE=./configs/server.json \
	LISTEN_ADDR=":8000" \
	 ./target/server

kick: target/kick
	CONF_FILE=./configs/server.json \
	target/kick

kick-setupai: kick
	CONF_FILE=./configs/server.json \
	target/kick setupai

kick-playout: kick
	CONF_FILE=./configs/server.json \
	target/kick playout

test:
	go test -v ./lib/...

build: target/server target/playout target/kick target/container


.PHONY: dev db test build kick
