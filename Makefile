.PHONY: build proto server client

build: client server

client:
	go build -o cli ./client/main.go

server:
	go build -o srv ./server/main.go

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/soda.proto
