package engine

import (
	"fmt"
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

type ThingNotFoundError struct {
	error

	Key string
}

func NewThingNotFoundError(key string) ThingNotFoundError {
	return ThingNotFoundError{
		Key: key,
	}
}

func (err ThingNotFoundError) Error() string {
	return fmt.Sprintf(`Not found: "%s"`, err.Key)
}

// func (err ThingNotFoundError) Error() string {}
