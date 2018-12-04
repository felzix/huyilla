package main

import (
	"github.com/gogo/protobuf/proto"
	"github.com/peterbourgon/diskv"
	"github.com/satori/go.uuid"
	"regexp"
	"strings"
)

func makeDB(saveDir string, cacheSize uint64) (*diskv.Diskv, error) {
	unique, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	return diskv.New(diskv.Options{
		BasePath:          saveDir,
		TempDir:           "/tmp/tempdir-huyilla-" + unique.String(),
		AdvancedTransform: filesystemTransform,
		InverseTransform:  filesystemInverseTransform,
		CacheSizeMax:      cacheSize,
	}), nil
}

func filesystemTransform(key string) *diskv.PathKey {
	path := strings.Split(key, ".")
	last := len(path) - 1
	return &diskv.PathKey{
		Path:     path[:last],
		FileName: path[last],
	}
}

func filesystemInverseTransform(pathKey *diskv.PathKey) (key string) {
	return strings.Join(pathKey.Path, "/") + pathKey.FileName
}

var regexFileNotFound = regexp.MustCompile("no such file or directory")

func fileIsNotFound(err error) bool {
	return regexFileNotFound.MatchString(err.Error())
}

func gettum(world *World, key string, thing proto.Unmarshaler) error {
	if blob, err := world.DB.Read(key); err == nil {
		if err := thing.Unmarshal(blob); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func settum(world *World, key string, thing proto.Marshaler) error {
	if blob, err := thing.Marshal(); err == nil {
		if err := world.DB.Write(key, blob); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func enddum(world *World, key string) error {
	return world.DB.Erase(key)
}

func hassum(world *World, key string) bool {
	return world.DB.Has(key)
}
