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

test: proto
	go test -gcflags=-l $(PKG)/...

deps:
	go get \
		github.com/gogo/protobuf/jsonpb \
		github.com/gogo/protobuf/proto \
		github.com/spf13/cobra \
		github.com/pkg/errors \

clean:
	go clean
	rm -f \
		types/types.pb.go \
		run/engine \
		run/textclient \
