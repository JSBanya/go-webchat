BINARY=server
export GOPATH=$(CURDIR)/vendor/:$(CURDIR)/cmd/

all: build

build:
	go fmt $(CURDIR)/cmd/*.go
	go build -o $(BINARY) $(CURDIR)/cmd/*.go

fetch:
	go get github.com/gorilla/websocket

clean:
	go clean
	rm -rf vendor/*
	rm -f $(BINARY)

.PHONY: all build fetch clean