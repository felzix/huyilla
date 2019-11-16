package types

import (
	"fmt"
)

type Database interface {
	Get(key string, thing Serializable) error
	Set(key string, thing Serializable) error
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
