.PHONY: help all client server clean

help:
	# all    - build all binaries
	# client - build client
	# server - build server
	# clean  - clean binaries

all: client server

client:
	go build -o client cmd/client/main.go

server:
	go build -o server cmd/server/main.go

clean:
	rm -rf client
	rm -rf server
