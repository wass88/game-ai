LIB = $(wildcard lib/*)

GO = go build
GO_ALPINE = GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

target/server: cmd/server/main.go $(LIB)
	$(GO) -o $@ cmd/server/main.go

target/playout: cmd/playout/main.go $(LIB)
	$(GO) -o $@ cmd/playout/main.go

target/server.alpine: cmd/server/main.go $(LIB)
	$(GO_ALPINE) -o $@ cmd/server/main.go

docker-dev: target/server.alpine
	docker-compose build --no-cache
	docker-compose up 

db:
	docker-compose up -d db

dev: target/server
	MYSQL_DATABASE="root:goodpassword@tcp(127.0.0.1:13306)/dev" ./target/server

.PHONY: dev db
