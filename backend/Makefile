GO_ALPINE=GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build
GO=go build

build: target/server target/playout target/kick target/container

target/server: cmd/server/main.go  $(wildcard lib/server/*) 
	$(GO) -o $@ cmd/server/main.go

target/playout: cmd/playout/main.go $(wildcard lib/playout/*) $(wildcard lib/game/*) $(wildcard lib/game/*/*)
	$(GO) -o $@ cmd/playout/main.go

target/container: cmd/container/main.go $(wildcard lib/container/*) $(wildcard lib/docker/*)
	$(GO) -o $@ cmd/container/main.go

target/kick: cmd/kick/main.go $(wildcard lib/server/*)
	$(GO) -o $@ cmd/kick/main.go

bin: bin/server bin/playout bin/container bin/kick

bin/server: cmd/server/main.go  $(wildcard lib/server/*) 
	$(GO_ALPINE) -o $@ cmd/server/main.go

bin/playout: cmd/playout/main.go $(wildcard lib/playout/*) $(wildcard lib/game/*) $(wildcard lib/game/*/*)
	$(GO_ALPINE) -o $@ cmd/playout/main.go

bin/container: cmd/container/main.go $(wildcard lib/container/*) $(wildcard lib/docker/*)
	$(GO_ALPINE) -o $@ cmd/container/main.go

bin/kick: cmd/kick/main.go $(wildcard lib/server/*)
	$(GO_ALPINE) -o $@ cmd/kick/main.go


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


.PHONY: dev db test build kick bin