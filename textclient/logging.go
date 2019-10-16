package main

import (
	"github.com/op/go-logging"
	"os"
)

var Log = logging.MustGetLogger("huyilla-textclient")

func MakeFileLogBackend(filename string) (logging.LeveledBackend, func() error, error) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, nil, err
	}

	backend := logging.NewLogBackend(f, "", 0)
	backendLeveled := logging.AddModuleLevel(backend)
	backendLeveled.SetLevel(logging.DEBUG, "")

	return backendLeveled, f.Close, nil
}
