package engine

import (
	"github.com/gogo/protobuf/proto"
)

type Database interface {
	Get(key string, thing proto.Unmarshaler) error
	Set(key string, thing proto.Marshaler) error
	GetByPrefix(prefix string) <-chan string
	Has(key string) bool
	End(key string) error
	EndAll() error
}
