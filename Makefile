PKG = github.com/felzix/huyilla

.PHONY: all clean test lint deps proto cli textclient

all: external-plugin

internal-plugin: huyilla.so
external-plugin: huyilla.0.0.1

huyilla.0.0.1: proto
	mkdir -p run/contracts
	go build -o run/contracts/$@ ./contract

cli: proto
	go build -o run/cli ./cli

textclient: proto
	go build -o run/textclient ./textclient

huyilla.so: proto
	mkdir -p run/contracts
	go build -buildmode=plugin -o run/contracts/$@ ./contract

%.pb.go: %.proto
	protoc --gofast_out=. $<

proto: types/types.pb.go

test: proto
	go test $(PKG)/...

lint:
	golint ./...

deps:
	go get \
		github.com/gogo/protobuf/jsonpb \
		github.com/gogo/protobuf/proto \
		github.com/spf13/cobra \
		github.com/gomodule/redigo/redis \
		github.com/gorilla/websocket \
		github.com/pkg/errors \
		github.com/grpc-ecosystem/go-grpc-prometheus \
        github.com/hashicorp/go-plugin \
		github.com/loomnetwork/go-loom

clean:
	go clean
	rm -f \
		types/types.pb.go \
		run/cli \
		run/textclient \
		run/contracts/huyilla.so \
		run/contracts/huyilla.0.0.1 \