.PHONY: help all client server clean

help:
	# all    - build all binaries
	# build  - build client and server
	# test   - test target files
	# clean  - clean binaries

all: test build

build:
	mkdir -p bin
	go build -o bin/client cmd/client/main.go
	go build -o bin/server cmd/server/main.go

test: go_test

clean:
	rm -rf client
	rm -rf server

go_test:
	go test $$(go list ./pkg/... ./cmd/....)
