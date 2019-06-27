PKG = github.com/felzix/huyilla

.PHONY: all clean test deps proto rundir engine textclient

all: engine textclient

engine: proto rundir
	go build -o run/engine ./engine

textclient: proto rundir
	go build -o run/textclient ./textclient

rundir:
	mkdir -p run

%.pb.go: %.proto
	protoc --gofast_out=. $<

proto: types/types.pb.go

fmt:
	go fmt $(PKG)/...


# Depends on engine to capture content's init() behavior.
test: engine
	go test -gcflags=-l $(PKG)/...

deps:
	go get \
		github.com/golang/protobuf/proto \
		github.com/gogo/protobuf/jsonpb \
		github.com/gogo/protobuf/proto \
		github.com/gogo/protobuf/protoc-gen-gofast \
		github.com/spf13/cobra \
		github.com/pkg/errors \
		github.com/dgrijalva/jwt-go \
		github.com/mitchellh/hashstructure \
		github.com/peterbourgon/diskv \
		github.com/satori/go.uuid \
		golang.org/x/crypto/bcrypt \
		github.com/grpc-ecosystem/go-grpc-prometheus \
		github.com/hashicorp/go-plugin \
		github.com/gorilla/mux \
		github.com/felzix/go-curses-react \
		github.com/felzix/goblin \


clean:
	go clean
	rm -f \
		types/types.pb.go \
		run/engine \
		run/textclient \
